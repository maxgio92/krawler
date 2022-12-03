package debian

import (
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
func (c *Debian) GetPackages(filter p.Filter) ([]p.Package, error) {
	// Merge custom config with default config.
	config, err := c.buildConfig(DebianDefaultConfig, c.config)
	if err != nil {
		return nil, err
	}

	// Build distribution version-specific mirror root URLs.
	releaseIndexURLs, err := c.buildReleaseIndexURLs(config.Mirrors, config.Versions)
	if err != nil {
		return nil, err
	}

	// TODO: from per-version (i.e. Debian Distribution) URL Release index file, get per-repository (e.g. Debian Component) Packages file
	// e.g. dists/stable/Release -> dists/stable/main/binary-amd64/Packages.xz

	// Build available repository URLs based on provided configuration,
	//for each distribution version.
	//
	// This should be part of Deb logic (i.e. deb.BuildPackagesIndexURLs function).
	packagesIndexURLs, err := c.buildPackagesIndexURLs(releaseIndexURLs, config.Repositories, c.vars)
	if err != nil {
		return nil, err
	}

	// TODO: from URLs of Packages index files, get deb Packages
	// e.g. dists/stable/main/binary-amd64/Packages.xz -> ... pool/main/l/linux-signed-amd64/linux-headers-amd64_5.10.140-1_amd64.deb

	// Get deb packages from each repository.
	//debs, err := deb.GetPackagesFromIndexURLs(packagesIndexURLs, filter.String(), filter.PackageFileNames()...)
	debs, err := deb.GetPackagesFromIndexURLs(packagesIndexURLs, filter.String(), filter.PackageFileNames()...)
	if err != nil {
		return nil, err
	}

	packages := make([]p.Package, len(debs))

	for i, v := range debs {
		v := v
		packages[i] = p.Package(&v)
	}

	return packages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Debian) buildReleaseIndexURLs(mirrors []p.Mirror, versions []distro.Version) ([]*url.URL, error) {
	versions, err := c.buildVersions(mirrors, versions)
	if err != nil {
		return []*url.URL{}, err
	}

	if (len(versions) > 0) && (len(mirrors) > 0) {
		var versionRoots []*url.URL

		for _, mirror := range mirrors {
			for _, version := range versions {
				v, err := url.JoinPath(mirror.URL, "dists", string(version))
				if err != nil {
					return nil, err
				}

				versionRoot, err := url.Parse(v)
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
func (c *Debian) buildVersions(mirrors []p.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
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
func (c *Debian) crawlVersions(mirrors []p.Mirror, debug bool) ([]distro.Version, error) {
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

	folderNames, err := scrape.CrawlFolders(seedUrls, DebianMirrorsDistroVersionRegex, false, debug)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

// Returns the list of repositories URLs.
func (c *Debian) buildPackagesIndexURLs(roots []*url.URL, components []p.Repository, vars map[string]interface{}) ([]*url.URL, error) {
	var urls []*url.URL

	paths, err := c.buildComponentPaths(components, vars)
	if err != nil {
		return []*url.URL{}, err
	}

	for _, root := range roots {
		//nolint:revive,stylecheck
		for _, path := range paths {
			// Get component URL from path.
			//nolint:revive,stylecheck
			us, err := url.JoinPath(root.String(), path)
			if err != nil {
				return nil, err
			}

			componentURL, err := url.Parse(us)
			if err != nil {
				return nil, err
			}

			urls = append(urls, componentURL)
		}
	}

	return urls, nil
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
