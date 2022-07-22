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

	perVersionMirrorUrls, err := c.buildPerVersionMirrorUrls(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	// Apply repository packages URI for each provided architecture.
	repositoriesUrls, err := c.buildRepositoriesUrls(perVersionMirrorUrls, config.Repositories, c.vars)
	if err != nil {
		return nil, err
	}

	packages, err = c.crawlPackages(repositoriesUrls, filter, debugScrape)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Centos) buildPerVersionMirrorUrls(mirrors []Mirror, versions []Version) ([]*url.URL, error) {
	versions, err := c.buildVersions(mirrors, versions)
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

// Returns the list of repositories URLs.
func (c *Centos) buildRepositoriesUrls(roots []*url.URL, repositories []Repository, vars map[string]interface{}) ([]*url.URL, error) {
	var urls []*url.URL

	uris, err := c.buildRepositoriesUris(repositories, vars)
	if err != nil {
		return []*url.URL{}, err
	}

	for _, root := range roots {
		//nolint:revive,stylecheck
		for _, uri := range uris {
			// Get repository URL from URI.
			//nolint:revive,stylecheck
			repositoryUrl, err := url.Parse(root.String() + uri)
			if err != nil {
				return nil, err
			}

			urls = append(urls, repositoryUrl)
		}
	}

	return urls, nil
}

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
func (c *Centos) crawlPackages(seedUrls []*url.URL, filter Filter, debug bool) ([]Package, error) {
	filenameRegex := `^` + string(filter) + `.+.` + CentosPackageFileExtension
	filenames, err := scrape.CrawlFiles(seedUrls, filenameRegex, debug)
	if err != nil {
		return []Package{}, err
	}
	var packages []Package
	for _, filename := range filenames {
		packages = append(packages, Package(filename))
	}

	return packages, nil
}
