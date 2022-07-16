package distro

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/maxgio92/krawler/pkg/template"
)

func init() {
	DistroByType[CentosType] = &Centos{}
}

type Centos struct {
	config Config
	vars   map[string]interface{}
}

func (c *Centos) Configure(config Config, vars map[string]interface{}) error {
	c.config = config
	c.vars = vars
	return nil
}

// For each mirror, for each distro version, for each repository,
// for each architecture, scrape.
func (c *Centos) GetPackages(filter Filter) ([]Package, error) {
	var packages []Package

	// Merge custom config with default config.
	config, err := c.buildConfig(CentosDefaultConfig, c.config)
	if err != nil {
		return nil, err
	}

	// Get distro versions for which to search packages.
	versions, err := c.buildVersions(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	versionsUrls, err := c.buildVersionsUrls(config.Mirrors, versions)
	if err != nil {
		return nil, err
	}

	// Apply repository packages URI for each provided architecture.
	//repositoriesUris, err := c.buildRepositoriesUris(config.Repositories, vars)
	repositoriesUris, err := c.buildRepositoriesUris(config.Repositories, c.vars)
	if err != nil {
		return nil, err
	}

	for _, root := range versionsUrls {
		//nolint:revive,stylecheck
		for _, repositoryUri := range repositoriesUris {
			// Get repository URL from URI.
			//nolint:revive,stylecheck
			repositoryUrl, err := url.Parse(root.String() + repositoryUri)
			if err != nil {
				return nil, err
			}

			// Crawl packages based on filter.
			p, err := c.crawlPackages(repositoryUrl, filter, debugScrape)
			if err != nil {
				return nil, err
			}

			packages = append(packages, p...)
		}
	}

	return packages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Centos) buildVersionsUrls(mirrors []Mirror, versions []Version) ([]*url.URL, error) {
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

	return nil, errNoDistroVersionSpecified
}

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (c *Centos) buildVersions(mirrors []Mirror, staticVersions []Version) ([]Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []Version

	dynamicVersions, err := c.crawlVersions(mirrors, debugScrape)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (c *Centos) crawlVersions(mirrors []Mirror, debug bool) ([]Version, error) {
	var versions []Version

	seedUrls := make([]*url.URL, 0, len(mirrors))
	for _, mirror := range mirrors {
		u, err := url.Parse(mirror.URL)
		if err != nil {
			return []Version{}, err
		}
		seedUrls = append(seedUrls, u)
	}

	folderNames, err := scrape.CrawlFolders(seedUrls, CentosMirrorsDistroVersionRegex, debug)
	if err != nil {
		return []Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, Version(v))
	}

	return versions, nil
}

// Returns the list of repositories URLs, for each specified architecture.
//
// TODO: the return of this method should compose the actual Scrape Config,
// which would consists of root URLs (of which below the URI segments)
// to scrape for pre-defined filters, i.e. package name.
func (c *Centos) buildRepositoriesUris(repositories []Repository, vars map[string]interface{}) ([]string, error) {
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
	if len(uris) < 1 {
		uris = append(uris, "/")
	}

	return uris, nil
}

// Returns the list of default repositories from the default config.
func (c *Centos) getDefaultRepositories() []Repository {
	var repositories []Repository

	for _, repository := range CentosDefaultConfig.Repositories {
		if !repositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}

// Returns a list of packages found on each page URL, filtered by filter.
//nolint:funlen,revive,stylecheck
func (c *Centos) crawlPackages(seedUrl *url.URL, filter Filter, debug bool) ([]Package, error) {

	filteredPackageRegex := `^` + string(filter) + `.+.` + CentosPackageFileExtension
	packageNames, err := scrape.CrawlFiles(seedUrl, filteredPackageRegex, debug)
	if err != nil {
		return []Package{}, err
	}
	var packages []Package
	for _, v := range packageNames {
		packages = append(packages, Package(v))
	}

	return packages, nil
}
