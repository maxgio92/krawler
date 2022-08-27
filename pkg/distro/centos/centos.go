package centos

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/maxgio92/krawler/pkg/template"
)

type Centos struct {
	config distro.Config
	vars   map[string]interface{}
}

func (c *Centos) Configure(config distro.Config, vars map[string]interface{}) error {
	c.config = config
	c.vars = vars
	return nil
}

// For each mirror, for each distro version, for each repository,
// for each architecture, get packages.
func (c *Centos) GetPackages(filter distro.Filter) ([]distro.Package, error) {
	var packages []distro.Package

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

	packages, err = rpm.GetPackagesFromRepositories(repositoriesUrls, string(filter), debugScrape)
	if err != nil {
		return nil, err
	}

	return packages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Centos) buildPerVersionMirrorUrls(mirrors []distro.Mirror, versions []distro.Version) ([]*url.URL, error) {
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

	return nil, distro.ErrNoDistroVersionSpecified
}

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (c *Centos) buildVersions(mirrors []distro.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := c.crawlVersions(mirrors, debugScrape)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (c *Centos) crawlVersions(mirrors []distro.Mirror, debug bool) ([]distro.Version, error) {
	var versions []distro.Version

	seedUrls := make([]*url.URL, 0, len(mirrors))
	for _, mirror := range mirrors {
		u, err := url.Parse(mirror.URL)
		if err != nil {
			return []distro.Version{}, err
		}
		seedUrls = append(seedUrls, u)
	}

	folderNames, err := scrape.CrawlFolders(seedUrls, CentosMirrorsDistroVersionRegex, debug)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

// Returns the list of repositories URLs.
func (c *Centos) buildRepositoriesUrls(roots []*url.URL, repositories []distro.Repository, vars map[string]interface{}) ([]*url.URL, error) {
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

func (c *Centos) buildRepositoriesUris(repositories []distro.Repository, vars map[string]interface{}) ([]string, error) {
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
func (c *Centos) getDefaultRepositories() []distro.Repository {
	var repositories []distro.Repository

	for _, repository := range CentosDefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}

// Returns a list of packages found on each page URL, filtered by filter.
func (c *Centos) crawlPackages(seedUrls []*url.URL, filter distro.Filter, debug bool) ([]distro.Package, error) {
	filenameRegex := `^` + string(filter) + `.+.` + CentosPackageFileExtension

	filenames, err := scrape.CrawlFiles(seedUrls, filenameRegex, debug)
	if err != nil {
		return []distro.Package{}, err
	}

	var packages []distro.Package
	for _, filename := range filenames {
		packages = append(packages, distro.Package(filename))
	}

	return packages, nil
}
