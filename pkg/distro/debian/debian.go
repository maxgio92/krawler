package debian

import (
	"net/url"
	"path"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/deb"
	"github.com/maxgio92/krawler/pkg/scrape"
)

type Debian struct {
	Config distro.Config
}

func (d *Debian) Configure(config distro.Config) error {
	c, err := d.BuildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}

	d.Config = c

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (d *Debian) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	d.Config.Output.Logger = options.Log()

	// Build distribution version-specific seed URLs.
	// TODO: introduce support for Release index files, where InRelease does not exist.
	distURLs, err := d.buildReleaseIndexURLs(d.Config.Mirrors, d.Config.Versions)
	if err != nil {
		return nil, err
	}

	components := []string{}
	for _, v := range d.Config.Repositories {
		components = append(components, strings.TrimPrefix(path.Clean(string(v.URI)), "/"))
	}

	searchOptions := deb.NewSearchOptions(&options, d.Config.Architectures, distURLs, components)

	debs, err := deb.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return debs, nil
}

// Returns the list of version-specific mirror URLs.
func (d *Debian) buildReleaseIndexURLs(mirrors []packages.Mirror, versions []distro.Version) ([]string, error) {
	versions, _ = d.buildVersions(mirrors, versions)

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
func (d *Debian) buildVersions(mirrors []packages.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := d.crawlVersions(mirrors)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (d *Debian) crawlVersions(mirrors []packages.Mirror) ([]distro.Version, error) {
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

	folderNames, err := scrape.CrawlFolders(
		seedUrls,
		DebianMirrorsDistroVersionRegex,
		false,
		d.Config.Output.Verbosity >= output.DebugLevel,
	)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

// Returns the list of default repositories from the default Config.
func (d *Debian) getDefaultRepositories() []packages.Repository {
	var repositories []packages.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
