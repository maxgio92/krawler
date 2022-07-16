package distro

import (
	"net/url"
	"strings"

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

func (c *Centos) buildConfig(def Config, user Config) (Config, error) {
	config, err := c.mergeConfig(def, user)
	if err != nil {
		return Config{}, err
	}

	err = c.sanitizeConfig(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// Returns the final configuration by merging the default with the user provided.
//nolint:unparam
func (c *Centos) mergeConfig(def Config, config Config) (Config, error) {
	if len(config.Archs) < 1 {
		config.Archs = def.Archs
	} else {
		for _, arch := range config.Archs {
			if arch == "" {
				config.Archs = def.Archs

				break
			}
		}
	}

	//nolint:nestif
	if len(config.Mirrors) < 1 {
		config.Mirrors = def.Mirrors
	} else {
		for _, mirror := range config.Mirrors {
			if mirror.URL == "" {
				config.Mirrors = def.Mirrors

				break
			}
		}
	}

	if len(config.Repositories) < 1 {
		config.Repositories = c.getDefaultRepositories()
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = c.getDefaultRepositories()

				break
			}
		}
	}

	return config, nil
}

// Returns the final configuration by overriding the default.
//nolint:unparam,unused
func (c *Centos) overrideConfig(def Config, override Config) (Config, error) {
	if len(override.Mirrors) > 0 {
		if override.Mirrors[0].URL != "" {
			return override, nil
		}
	}

	return def, nil
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

func (c *Centos) sanitizeConfig(config *Config) error {
	err := c.sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func (c *Centos) sanitizeMirrors(mirrors *[]Mirror) error {
	for i, mirror := range *mirrors {
		if !strings.HasSuffix(mirror.URL, "/") {
			(*mirrors)[i].URL = mirror.URL + "/"
		}

		_, err := url.Parse(mirror.URL)
		if err != nil {
			return err
		}
	}
	return nil
}
