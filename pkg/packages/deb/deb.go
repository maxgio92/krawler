package deb

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"path"
	"pault.ag/go/archive"
	"pault.ag/go/debian/deb"
	"strings"
)

func init() {
	log.SetLevel(log.FatalLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

// SearchPackages returns a slice of pault.ag/go/archive.Package objects, filtering as for search options.
// The function crawls the repositories with asynchronous and parallel workers.
func SearchPackages(so *SearchOptions) ([]archive.Package, error) {
	packages := []archive.Package{}

	// Run producer workers.
	for _, v := range so.SeedURLs() {
		distURL := v
		go searchPackagesFromDist(
			func() {
				so.Progress(1)
				so.WaitGroup().Done()
			},
			so, distURL)
	}

	// Run consumer worker.
	go consumePackages(so, &packages)

	// Wait for producers and consumers to complete.
	so.WaitAll()

	return packages, nil
}

// searchPackagesFromDist writes to a channel pault.ag/go/archive.Package objects, writes errors to a channel, through usage
// of asynchronous workers. It needs a function doneFunc to be executed on completion.
// Accepts as argument for filtering packages the package name as string and the deb dist URL where to look for packages.
func searchPackagesFromDist(doneFunc func(), so *SearchOptions, distURL string) {
	defer doneFunc()

	inRelease, err := getInReleaseFromDistURL(distURL)
	if err != nil {
		so.ErrCh() <- err
		return
	}

	indexURLs, err := getPackagesIndexURLsFromInRelease(inRelease, distURL)
	if err != nil {
		so.ErrCh() <- err
		return
	}

	indexSearchOptions := NewSearchOptions(so.PackageName(), indexURLs, fmt.Sprintf("Indexing packages for dist %s", path.Base(distURL)))

	// Run producer workers, to search packages from Packages index files.
	for _, v := range indexSearchOptions.SeedURLs() {
		if ExcludeInstallers && strings.Contains(v, "debian-installer") {
			indexSearchOptions.WaitGroup().Done()
			continue
		}

		go searchPackagesFromIndex(
			func() {
				indexSearchOptions.Progress(1)
				indexSearchOptions.WaitGroup().Done()
			},
			indexSearchOptions, v)
	}

	// Run consumer worker from child option set, to fill the parent search option set.
	go consumeIndexPackages(indexSearchOptions, so)

	// Wait for producers to complete.
	indexSearchOptions.WaitGroup().Wait()
	close(indexSearchOptions.ResultCh())
	close(indexSearchOptions.ErrCh())

	// Wait for consumers to complete.
	<-indexSearchOptions.DoneCh()
}

func consumePackages(so *SearchOptions, target *[]archive.Package) {
	resultCh := so.ResultCh()
	errCh := so.ErrCh()
	for errCh != nil || resultCh != nil {
		select {
		case p, ok := <-resultCh:
			if ok {
				log.Debug("Scanned DB")
				if len(p) > 0 {
					*target = append(*target, p...)
					log.Infof("New %d packages found", len(p))
				}
				continue
			}
			resultCh = nil
		case e, ok := <-errCh:
			if ok {
				log.Error(e)
				continue
			}
			errCh = nil
		}
	}
	so.DoneCh() <- true
}

// searchPackagesFromIndex searches and fills with a channel of deb packages from Packages index files.
// E.g. /dists/stable/main/binary-amd64/Packages.xz -> /pool/main/l/linux-signed-amd64/linux-headers-amd64_5.10.140-1_amd64.deb
func searchPackagesFromIndex(doneFunc func(), so *SearchOptions, indexURL string) {
	defer doneFunc()

	log.WithField("URL", indexURL).Debug("Downloading compressed index file")

	resp, err := http.Get(indexURL)
	if err != nil {
		so.ErrCh() <- err
		return
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		so.ErrCh() <- fmt.Errorf("download(%s): unexpected HTTP status code: got %d, want %d", indexURL, got, want)
		return
	}
	defer resp.Body.Close()

	log.WithField("URL", indexURL).Debug("Decompressing index file")

	debDecompressor := deb.DecompressorFor(PackagesIndexFormat)
	rd, err := debDecompressor(resp.Body)
	defer rd.Close()
	if err != nil {
		so.ErrCh() <- err
		return
	}

	log.WithField("URL", indexURL).Debug("Loading packages DB from index file")

	db, err := archive.LoadPackages(rd)
	if err != nil {
		so.ErrCh() <- err
		return
	}

	log.WithField("URL", indexURL).Debug("Querying packages from DB")

	query := func(p *archive.Package) bool {
		if strings.Contains(p.Package, so.PackageName()) && p.Architecture.CPU != "all" {
			return true
		}
		return false
	}

	p, err := db.Map(query)
	if err != nil {
		so.ErrCh() <- err
		return
	}

	so.ResultCh() <- p
}

func consumeIndexPackages(indexSearchOptions *SearchOptions, distSearchOptions *SearchOptions) {
	indexResultCh := indexSearchOptions.ResultCh()
	indexErrCh := indexSearchOptions.ErrCh()

	for indexErrCh != nil || indexResultCh != nil {
		select {
		case p, ok := <-indexResultCh:
			if ok {
				log.Debug("got a response from DB")
				if len(p) > 0 {
					distSearchOptions.ResultCh() <- p
				}
				continue
			}
			indexResultCh = nil
		case e, ok := <-indexErrCh:
			if ok {
				log.Debug("got an error from DB")
				distSearchOptions.ErrCh() <- e
				continue
			}
			indexErrCh = nil
		}
	}
	log.Debug("consumers are done")
	indexSearchOptions.DoneCh() <- true
}

// getInReleaseFromDistURL returns a *archive.Release object from the deb dist URL.
// It leverages pault.ag/go/archive and pault.ag/go/debian/deb libraries to parse and build the Release object.
func getInReleaseFromDistURL(distURL string) (*archive.Release, error) {
	inReleaseURL, err := url.JoinPath(distURL, InRelease)
	if err != nil {
		return nil, err
	}

	inReleaseResp, err := http.Get(inReleaseURL)
	if err != nil {
		return nil, err
	}
	if got, want := inReleaseResp.StatusCode, http.StatusOK; got != want {
		if inReleaseResp.StatusCode == 404 {
			return nil, fmt.Errorf("InRelease file not found with dist URL %s", distURL)
		}
		if inReleaseResp.StatusCode >= 500 && inReleaseResp.StatusCode < 600 {
			return nil, fmt.Errorf("internal error from mirror for release file with dist URL %s", distURL)
		}

		return nil, fmt.Errorf("download(%s): unexpected HTTP status code: got %d, want %d", inReleaseURL, got, want)
	}

	release, err := archive.LoadInRelease(inReleaseResp.Body, nil)
	if err != nil {
		return nil, err
	}

	return release, nil
}

// getPackagesIndexURLsFromInRelease returns from per dist Release index file, the URLs of the per component Packages
// index files.
// E.g. from /dists/stable/Release -> /dists/stable/main/binary-amd64/Packages.xz
func getPackagesIndexURLsFromInRelease(inRelease *archive.Release, distURL string) ([]string, error) {
	indexURLs := []string{}
	for _, v := range inRelease.MD5Sum {
		if strings.Contains(v.Filename, "Packages"+PackagesIndexFormat) {

			u, err := url.JoinPath(distURL, v.Filename)
			if err != nil {
				return nil, err
			}

			indexURLs = append(indexURLs, u)
		}
	}

	return indexURLs, nil
}
