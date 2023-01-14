package deb

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	progressbar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"pault.ag/go/archive"
	"pault.ag/go/debian/deb"
)

func init() {
	log.SetLevel(log.FatalLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
	})
}

// GetPackages returns a slice of pault.ag/go/archive.Package objects, filtering per package name.
// The needed arguments are the package name as string, and a list of URL of the deb dists as string
// slice.
func GetPackages(options *SearchOptions) ([]archive.Package, error) {
	packages := []archive.Package{}

	perDistWG := sync.WaitGroup{}
	perDistWG.Add(len(options.DistURLs()))

	packagesCh := make(chan []archive.Package)

	errCh := make(chan error)

	done := make(chan bool, 1)

	bar := progressbar.Default(int64(len(options.DistURLs())), "Total")

	// Run producer workers.
	for _, v := range options.DistURLs() {
		distURL := v
		go getDistPackages(bar, &perDistWG, packagesCh, errCh, options.PackageName(), distURL)
	}

	// Run consumer worker.
	go consumePackages(done, &packages, packagesCh, errCh)

	// Wait for producers to complete.
	perDistWG.Wait()
	close(packagesCh)
	close(errCh)

	// Wait for consumers to complete.
	<-done

	return packages, nil
}

// getDistPackages writes to a channel pault.ag/go/archive.Package objects, writes errors to a channel, through usage
// of asynchronous workers. It needs a *sync.WaitGroup as arguments to synchronize the producer workers.
// Accepts as argument for filtering packages the package name as string and the deb dist URL where to look for packages.
func getDistPackages(bar *progressbar.ProgressBar, waitGroup *sync.WaitGroup, packagesCh chan []archive.Package, errCh chan error, packageName string, distURL string) {
	defer waitGroup.Done()
	defer bar.Add(1)

	indexesWG := sync.WaitGroup{}

	packagesInternalCh := make(chan []archive.Package)

	errInternalCh := make(chan error)

	done := make(chan bool, 1)

	inRelease, err := getInReleaseFromDistURL(distURL)
	if err != nil {
		errCh <- err
		return
	}

	// From per-dist Release index file, get per-component Packages index file.
	// E.g. /dists/stable/Release -> /dists/stable/main/binary-amd64/Packages.xz
	indexPaths := []string{}
	for _, v := range inRelease.MD5Sum {
		if strings.Contains(v.Filename, "Packages"+PackagesIndexFormat) {
			indexPaths = append(indexPaths, v.Filename)
		}
	}

	indexesWG.Add(len(indexPaths))
	internalBar := progressbar.Default(int64(len(indexPaths)), fmt.Sprintf("Indexing packages for dist %s", path.Base(distURL)))

	// From Packages index files, get deb Packages.
	// E.g. /dists/stable/main/binary-amd64/Packages.xz -> /pool/main/l/linux-signed-amd64/linux-headers-amd64_5.10.140-1_amd64.deb
	//
	// Run producer workers.
	for _, v := range indexPaths {
		if ExcludeInstallers && strings.Contains(v, "debian-installer") {
			indexesWG.Done()
			continue
		}

		indexURL, err := url.JoinPath(distURL, v)
		if err != nil {
			errCh <- err
			indexesWG.Done()
			return
		}

		go getIndexPackages(internalBar, &indexesWG, packagesInternalCh, errInternalCh, packageName, indexURL)
	}

	// Run consumer worker.
	go func() {
		for errInternalCh != nil || packagesInternalCh != nil {
			select {
			case p, ok := <-packagesInternalCh:
				if ok {
					log.Debug("got a response from DB")
					if len(p) > 0 {
						packagesCh <- p
					}
					continue
				}
				packagesInternalCh = nil
			case e, ok := <-errInternalCh:
				if ok {
					log.Debug("got an error from DB")
					errCh <- e
					continue
				}
				errInternalCh = nil
			}
		}
		log.Debug("consumers are done")
		done <- true
	}()

	// Wait for producers to complete.
	indexesWG.Wait()
	close(packagesInternalCh)
	close(errInternalCh)

	// Wait for consumers to complete.
	<-done
}

func consumePackages(done chan bool, packages *[]archive.Package, packagesCh chan []archive.Package, errCh chan error) {
	for errCh != nil || packagesCh != nil {
		select {
		case p, ok := <-packagesCh:
			if ok {
				log.Debug("Scanned DB")
				if len(p) > 0 {
					*packages = append(*packages, p...)
					log.Infof("New %d packages found", len(p))
				}
				continue
			}
			packagesCh = nil
		case e, ok := <-errCh:
			if ok {
				log.Error(e)
				continue
			}
			errCh = nil
		}
	}
	done <- true
}

func getIndexPackages(bar *progressbar.ProgressBar, waitGroup *sync.WaitGroup, packagesCh chan []archive.Package, errCh chan error, packageName string, indexURL string) {
	defer waitGroup.Done()
	defer bar.Add(1)

	log.WithField("URL", indexURL).Debug("Downloading compressed index file")

	resp, err := http.Get(indexURL)
	if err != nil {
		errCh <- err
		return
	}
	if got, want := resp.StatusCode, http.StatusOK; got != want {
		errCh <- fmt.Errorf("download(%s): unexpected HTTP status code: got %d, want %d", indexURL, got, want)
		return
	}
	defer resp.Body.Close()

	log.WithField("URL", indexURL).Debug("Decompressing index file")

	debDecompressor := deb.DecompressorFor(PackagesIndexFormat)
	rd, err := debDecompressor(resp.Body)
	defer rd.Close()
	if err != nil {
		errCh <- err
		return
	}

	log.WithField("URL", indexURL).Debug("Loading packages DB from index file")

	db, err := archive.LoadPackages(rd)
	if err != nil {
		errCh <- err
		return
	}

	log.WithField("URL", indexURL).Debug("Querying packages from DB")

	query := func(p *archive.Package) bool {
		if strings.Contains(p.Package, packageName) && p.Architecture.CPU != "all" {
			return true
		}
		return false
	}

	p, err := db.Map(query)
	if err != nil {
		errCh <- err
		return
	}

	packagesCh <- p
}

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
