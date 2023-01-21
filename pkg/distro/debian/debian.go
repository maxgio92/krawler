package debian

import (
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages/deb"
	"net/url"

	"github.com/maxgio92/krawler/pkg/distro"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/maxgio92/krawler/pkg/utils/template"
)

type Debian struct {
	config distro.Config
	vars   map[string]interface{}
}

func (c *Debian) Configure(config distro.Config, vars map[string]interface{}) error {
	c.config = config
	c.vars = vars

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (c *Debian) GetPackages(options p.PackageOptions) ([]p.Package, error) {
	// Merge custom config with default config.
	config, err := c.buildConfig(DebianDefaultConfig, c.config)
	if err != nil {
		return nil, err
	}

	// Build distribution version-specific mirror root URLs.
	// TODO: introduce support for Release index files, where InRelease does not exist.
	distURLs, err := c.buildReleaseIndexURLs(config.Mirrors, config.Versions, options.Verbosity())
	if err != nil {
		return nil, err
	}

	searchOptions := deb.NewSearchOptions(options.PackageName(), distURLs, options.Verbosity())

	debs, err := deb.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	packages := make([]p.Package, len(debs))

	for i, v := range debs {
		v := v
		packages[i] = &deb.Package{
			Name:    v.Package,
			Arch:    v.Architecture.String(),
			Version: v.Version.String(),
		}
	}

	return packages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Debian) buildReleaseIndexURLs(mirrors []p.Mirror, versions []distro.Version, verbosity output.Verbosity) ([]string, error) {
	versions, err := c.buildVersions(mirrors, versions, verbosity)
	if err != nil {
		return nil, nil
	}

	if (len(versions) > 0) && (len(mirrors) > 0) {
		var versionRoots []string

		for _, mirror := range mirrors {
			for _, version := range versions {
				v, err := url.JoinPath(mirror.URL, "dists", string(version))
				if err != nil {
					return nil, err
				}

				versionRoots = append(versionRoots, v)
			}
		}

		return versionRoots, nil
	}

	return nil, distro.ErrNoDistroVersionSpecified
}

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (c *Debian) buildVersions(mirrors []p.Mirror, staticVersions []distro.Version, verbosity output.Verbosity) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := c.crawlVersions(mirrors, verbosity)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (c *Debian) crawlVersions(mirrors []p.Mirror, verbosity output.Verbosity) ([]distro.Version, error) {
	versions := []distro.Version{}

	seedUrls := make([]*url.URL, 0, len(mirrors))

	for _, mirror := range mirrors {
		distsURL, err := url.JoinPath(mirror.URL, "dists/")
		if err != nil {
			return []distro.Version{}, err
		}

		u, err := url.Parse(distsURL)
		if err != nil {
			return []distro.Version{}, err
		}

		seedUrls = append(seedUrls, u)
	}

	debug := false
	if verbosity >= output.DebugLevel {
		debug = true
	}

	folderNames, err := scrape.CrawlFolders(seedUrls, DebianMirrorsDistroVersionRegex, false, debug)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

func (c *Debian) buildComponentPaths(components []p.Repository, vars map[string]interface{}) ([]string, error) {
	paths := []string{}

	for _, component := range components {
		if component.URI != "" {
			result, err := template.MultiplexAndExecute(string(component.URI), vars)
			if err != nil {
				return nil, err
			}

			paths = append(paths, result...)
		}
	}

	// Scrape for all possible components.
	if len(paths) < 1 {
		paths = append(paths, "/")
	}

	return paths, nil
}

// Returns the list of default repositories from the default config.
func (c *Debian) getDefaultRepositories() []p.Repository {
	var repositories []p.Repository

	for _, repository := range DebianDefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
