package amazonlinux2022

import (
	"github.com/maxgio92/krawler/pkg/distro"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/maxgio92/krawler/pkg/utils/template"
	"net/url"
)

type AmazonLinux2022 struct {
	config distro.Config
	vars   map[string]interface{}
}

func (a *AmazonLinux2022) Configure(config distro.Config, vars map[string]interface{}) error {
	a.config = config
	a.vars = vars

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (a *AmazonLinux2022) GetPackages(filter p.Filter) ([]p.Package, error) {
	config, err := a.buildConfig(DefaultConfig, a.config)
	if err != nil {
		return nil, err
	}

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorURLs, err := a.buildMirrorURLs(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	//for each distribution version.
	repositoriesURLs, err := a.buildRepositoriesURLs(perVersionMirrorURLs, config.Repositories, a.vars)
	if err != nil {
		return nil, err
	}

	// Dereference repository URLs.
	results := []*url.URL{}

	for _, ar := range config.Archs {
		for _, v := range repositoriesURLs {
			r, err := a.dereferenceRepositoryURL(v, ar)
			if err != nil {
				return nil, err
			}

			results = append(results, r)
		}
	}

	// Get RPM packages from each repository.
	rpmPackages, err := rpm.GetPackagesFromRepositories(results, filter.String(), filter.PackageFileNames()...)
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

// Returns the list of version-specific mirror URLs.
func (a *AmazonLinux2022) buildMirrorURLs(mirrors []p.Mirror, versions []distro.Version) ([]*url.URL, error) {
	versions, err := a.buildVersions(mirrors, versions)
	if err != nil {
		return []*url.URL{}, err
	}

	if (len(versions) > 0) && (len(mirrors) > 0) {
		var versionRoots []*url.URL

		for _, mirror := range mirrors {
			for _, version := range versions {
				versionRoot, err := url.Parse(mirror.URL + string(version))
				if err != nil {
					return nil, err
				}

				versionRoots = append(versionRoots, versionRoot)
			}
		}

		return versionRoots, nil
	}

	return nil, distro.ErrNoDistroVersionSpecified
}

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (a *AmazonLinux2022) buildVersions(mirrors []p.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := a.crawlVersions(mirrors, debugScrape)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (a *AmazonLinux2022) crawlVersions(mirrors []p.Mirror, debug bool) ([]distro.Version, error) {
	versions := []distro.Version{}

	seedUrls := make([]*url.URL, 0, len(mirrors))

	for _, mirror := range mirrors {
		u, err := url.Parse(mirror.URL)
		if err != nil {
			return []distro.Version{}, err
		}

		seedUrls = append(seedUrls, u)
	}

	folderNames, err := scrape.CrawlFolders(seedUrls, MirrorsDistroVersionRegex, debug)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

// Returns the list of repositories URLs.
func (a *AmazonLinux2022) buildRepositoriesURLs(roots []*url.URL, repositories []p.Repository, vars map[string]interface{}) ([]*url.URL, error) {
	var urls []*url.URL

	uris, err := a.buildRepositoriesURIs(repositories, vars)
	if err != nil {
		return []*url.URL{}, err
	}

	for _, root := range roots {
		for _, uri := range uris {

			us, err := url.JoinPath(root.String(), uri)
			if err != nil {
				return nil, err
			}

			repositoryUrl, err := url.Parse(us)
			if err != nil {
				return nil, err
			}

			urls = append(urls, repositoryUrl)
		}
	}

	return urls, nil
}

func (a *AmazonLinux2022) buildRepositoriesURIs(repositories []p.Repository, vars map[string]interface{}) ([]string, error) {
	uris := []string{}

	for _, repository := range repositories {
		if repository.URI != "" {
			result, err := template.MultiplexAndExecute(string(repository.URI), vars)
			if err != nil {
				return nil, err
			}

			uris = append(uris, result...)
		}
	}

	// Scrape for all possible repositories.
	if len(uris) == 0 {
		uris = append(uris, "/")
	}

	return uris, nil
}
