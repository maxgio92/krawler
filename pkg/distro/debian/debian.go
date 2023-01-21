package debian

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/deb"
	"github.com/maxgio92/krawler/pkg/scrape"
)

type Debian struct {
	config distro.Config
}

func (d *Debian) Configure(config distro.Config) error {
	c, err := d.buildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}
	d.config = c

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (d *Debian) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	d.config.Output.Logger = options.Log()

	// Build distribution version-specific seed URLs.
	// TODO: introduce support for Release index files, where InRelease does not exist.
	distURLs, err := d.buildReleaseIndexURLs(d.config.Mirrors, d.config.Versions, options.Verbosity())
	if err != nil {
		return nil, err
	}

	searchOptions := deb.NewSearchOptions(&options, d.config.Architectures, distURLs)

	debs, err := deb.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return debs, nil
}

// Returns the list of version-specific mirror URLs.
func (d *Debian) buildReleaseIndexURLs(mirrors []packages.Mirror, versions []distro.Version, verbosity output.Verbosity) ([]string, error) {
	versions, err := d.buildVersions(mirrors, versions, verbosity)
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
func (d *Debian) buildVersions(mirrors []packages.Mirror, staticVersions []distro.Version, verbosity output.Verbosity) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := d.crawlVersions(mirrors, verbosity)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (d *Debian) crawlVersions(mirrors []packages.Mirror, verbosity output.Verbosity) ([]distro.Version, error) {
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

// Returns the list of default repositories from the default config.
func (d *Debian) getDefaultRepositories() []packages.Repository {
	var repositories []packages.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
