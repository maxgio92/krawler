package v2

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	common "github.com/maxgio92/krawler/pkg/distro/amazonlinux"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
)

type AmazonLinux struct {
	common.AmazonLinux
}

func (a *AmazonLinux) Configure(config distro.Config) error {
	return a.ConfigureCommon(DefaultConfig, config)
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (a *AmazonLinux) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	a.Config.Output.Logger = options.Log()

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorURLs, err := a.BuildMirrorURLs(a.Config.Mirrors, a.Config.Versions)
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	repositoriesURLrefs, err := common.BuildRepositoriesURLs(perVersionMirrorURLs, a.Config.Repositories)
	if err != nil {
		return nil, err
	}

	// Dereference repository URLs.
	repositoryURLs, err := a.dereferenceRepositoryURLs(repositoriesURLrefs, a.Config.Archs)
	if err != nil {
		return nil, err
	}

	// Get RPM packages from each repository.
	rss := []string{}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}

	searchOptions := rpm.NewSearchOptions(&options, a.Config.Archs, rss)
	rpmPackages, err := rpm.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return rpmPackages, nil
}

func (a *AmazonLinux) dereferenceRepositoryURLs(repoURLs []*url.URL, archs []packages.Architecture) ([]*url.URL, error) {
	var urls []*url.URL

	for _, ar := range archs {
		for _, v := range repoURLs {
			r, err := a.dereferenceRepositoryURL(v, ar)
			if err != nil {
				return nil, err
			}

			if r != nil {
				urls = append(urls, r)
			}
		}
	}

	return urls, nil
}

func (a *AmazonLinux) dereferenceRepositoryURL(src *url.URL, arch packages.Architecture) (*url.URL, error) {
	var dest *url.URL

	mirrorListURL, err := url.JoinPath(src.String(), string(arch), "mirror.list")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, mirrorListURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.Config.Output.Logger.Error("Amazon Linux v2 repository URL not valid to be dereferenced")
		//nolint:nilnil
		return nil, nil
	}

	if resp.Body == nil {
		a.Config.Output.Logger.Error("empty response from Amazon Linux v2 repository reference URL")
		//nolint:nilnil
		return nil, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Get first repository URL available, no matter what the geolocation.
	s := strings.Split(string(b), "\n")[0]

	dest, err = url.Parse(s)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
