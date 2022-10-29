package v2022

import (
	"context"
	"fmt"
	"github.com/maxgio92/krawler/pkg/distro"
	common "github.com/maxgio92/krawler/pkg/distro/amazonlinux"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type AmazonLinux struct {
	common.AmazonLinux
}

func (a *AmazonLinux) Configure(config distro.Config, vars map[string]interface{}) error {
	return a.ConfigureCommon(config, vars)
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (a *AmazonLinux) GetPackages(filter p.Filter) ([]p.Package, error) {
	config, err := common.BuildConfig(DefaultConfig, a.Config)
	if err != nil {
		return nil, err
	}

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorURLs, err := common.BuildMirrorURLs(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	//for each distribution version.
	repositoriesURLrefs, err := common.BuildRepositoriesURLs(perVersionMirrorURLs, config.Repositories, a.Vars)
	if err != nil {
		return nil, err
	}

	// Dereference repository URLs.
	repositoryURLs, err := dereferenceRepositoryURLs(repositoriesURLrefs, config.Archs)
	if err != nil {
		return nil, err
	}

	// Get RPM packages from each repository.
	rpmPackages, err := rpm.GetPackagesFromRepositories(repositoryURLs, filter.String(), filter.PackageFileNames()...)
	if err != nil {
		return nil, err
	}

	packages := make([]p.Package, len(rpmPackages))

	for i, v := range rpmPackages {
		v := v
		packages[i] = p.Package(&v)
	}

	return packages, nil
}

func dereferenceRepositoryURLs(repoURLs []*url.URL, archs []distro.Arch) ([]*url.URL, error) {
	var urls []*url.URL

	for _, ar := range archs {
		for _, v := range repoURLs {
			r, err := dereferenceRepositoryURL(v, ar)
			if err != nil {
				return nil, err
			}

			urls = append(urls, r)
		}
	}

	return urls, nil
}

func dereferenceRepositoryURL(src *url.URL, arch distro.Arch) (*url.URL, error) {
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
		return nil, fmt.Errorf("Amazon Linux v2 repository URL not valid to be dereferenced")
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response from Amazon Linux v2 repository reference URL")
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
