package deb

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/maxgio92/krawler/pkg/packages"

	"golang.org/x/exp/slices"
	"pault.ag/go/archive"
	"pault.ag/go/debian/deb"
)

// SearchPackages returns a slice of pault.ag/go/archive.Package objects, filtering as for search options.
// The function crawls the repositories with asynchronous and parallel workers.
func SearchPackages(so *SearchOptions) ([]packages.Package, error) {
	var result []packages.Package

	search := func(distURL string) {
		searchPackagesFromDist(
			func() {
				so.Progress(1)
				so.SigProducerCompletion()
			},
			so, distURL)
	}

	collect := func() {
		so.Consume(
			func(p ...packages.Package) {
				so.Log().Debug("Scanned DB")
				if len(p) > 0 {
					result = append(result, p...)
					so.Log().Infof("New %d packages found", len(p))
				}
			},
			func(e error) {
				so.Log().Error(e)
			},
		)
	}

	// Run search producers.
	for _, v := range so.SeedURLs() {
		distURL := v
		go search(distURL)
	}

	// Run collect consumer.
	go collect()

	// Wait for producers and consumers to complete and cleanup.
	so.WaitAndClose()

	return result, nil
}

// searchPackagesFromDist writes to a channel pault.ag/go/archive.Package objects, writes errors to a channel, through usage
// of asynchronous workers. It needs a function doneFunc to be executed on completion.
// Accepts as argument for filtering packages the package name as string and the deb dist URL where to look for packages.
func searchPackagesFromDist(doneFunc func(), distSO *SearchOptions, distURL string) {
	defer doneFunc()

	inRelease, err := getInReleaseFromDistURL(distURL)
	if err != nil {
		distSO.SendError(err)

		return
	}

	indexURLs, err := getPackagesIndexURLsFromInRelease(inRelease, distURL)
	if err != nil {
		distSO.SendError(err)

		return
	}

	o := packages.NewSearchOptions(distSO.PackageName(), distSO.Architectures(), indexURLs, distSO.Verbosity(), fmt.Sprintf("Indexing packages for dist %s", path.Base(distURL)))
	indexSO := NewSearchOptions(o, o.Architectures(), o.SeedURLs(), distSO.Components())

	// Run producers, to search packages from Packages index files.
	for _, v := range indexSO.SeedURLs() {
		if ExcludeInstallers && strings.Contains(v, "debian-installer") {
			indexSO.SigProducerCompletion()

			continue
		}

		ss := strings.Split(v, "/")
		component := ss[len(ss)-3]

		if !slices.Contains(indexSO.Components(), component) {
			indexSO.SigProducerCompletion()

			continue
		}

		go searchPackagesFromIndex(
			func() {
				indexSO.Progress(1)
				indexSO.SigProducerCompletion()
			},
			indexSO, v)
	}

	// Run consumer from child option set, to fill the parent search option set.
	go indexSO.Consume(
		func(p ...packages.Package) {
			indexSO.Log().Debug("got a response from DB")
			if len(p) > 0 {
				distSO.SendMessage(p...)
			}
		},
		func(e error) {
			indexSO.Log().Debug("got an error from DB")
			distSO.SendError(e)
		},
	)

	// Wait for producersWG and consumer to complete.
	indexSO.WaitAndClose()
}

// searchPackagesFromIndex searches and fills with a channel of deb packages from Packages index files.
// E.g. /dists/stable/main/binary-amd64/Packages.xz -> /pool/main/l/linux-signed-amd64/linux-headers-amd64_5.10.140-1_amd64.deb
//
//nolint:funlen,cyclop
func searchPackagesFromIndex(doneFunc func(), so *SearchOptions, indexURL string) {
	defer doneFunc()

	so.Log().WithField("URL", indexURL).Debug("Downloading compressed index file")

	//nolint:gosec,noctx
	resp, err := http.Get(indexURL)
	if err != nil {
		so.SendError(err)

		return
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		//nolint:goerr113
		so.SendError(fmt.Errorf("download(%s): unexpected HTTP status code: got %d, want %d", indexURL, got, want))

		return
	}
	defer resp.Body.Close()

	so.Log().WithField("URL", indexURL).Debug("Decompressing index file")

	debDecompressor := deb.DecompressorFor(PackagesIndexFormat)
	rd, err := debDecompressor(resp.Body)

	defer rd.Close()

	if err != nil {
		so.SendError(err)

		return
	}

	so.Log().WithField("URL", indexURL).Debug("Loading packages DB from index file")

	db, err := archive.LoadPackages(rd)
	if err != nil {
		so.SendError(err)

		return
	}

	so.Log().WithField("URL", indexURL).Debug("Querying packages from DB")

	query := func(p *archive.Package) bool {
		if strings.Contains(p.Package, so.PackageName()) {
			if p.Architecture.CPU == "all" {
				return false
			}

			if slices.Contains(so.Architectures(), packages.Architecture(p.Architecture.CPU)) {
				return true
			}

			if so.Architectures() == nil {
				return true
			}
		}

		return false
	}

	ds, err := db.Map(query)
	if err != nil {
		so.SendError(err)

		return
	}

	// Convert deb packages to a standard type.
	ps := []packages.Package{}

	for _, d := range ds {
		p := &Package{
			Name:    d.Package,
			Arch:    d.Architecture.String(),
			Version: d.Version.String(),
		}
		ps = append(ps, p)
	}

	so.SendMessage(ps...)
}

// getInReleaseFromDistURL returns a *archive.Release object from the deb dist URL.
// It leverages pault.ag/go/archive and pault.ag/go/debian/deb libraries to parse and build the Release object.
func getInReleaseFromDistURL(distURL string) (*archive.Release, error) {
	inReleaseURL, err := url.JoinPath(distURL, InRelease)
	if err != nil {
		return nil, err
	}

	//nolint:gosec,noctx
	inReleaseResp, err := http.Get(inReleaseURL)
	if err != nil {
		return nil, err
	}
	defer inReleaseResp.Body.Close()

	if got, want := inReleaseResp.StatusCode, http.StatusOK; got != want {
		if inReleaseResp.StatusCode == http.StatusNotFound {
			//nolint:goerr113
			return nil, fmt.Errorf("InRelease file not found with dist URL %s", distURL)
		}

		if inReleaseResp.StatusCode >= 500 && inReleaseResp.StatusCode < 600 {
			//nolint:goerr113
			return nil, fmt.Errorf("internal error from mirror for release file with dist URL %s", distURL)
		}

		//nolint:goerr113
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
// E.g. from /dists/stable/Release -> /dists/stable/main/binary-amd64/Packages.xz.
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
