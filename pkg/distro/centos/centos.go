package centos

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"github.com/maxgio92/krawler/pkg/scrape"
)

type Centos struct {
	config distro.Config
}

func (c *Centos) Configure(config distro.Config) error {
	cfg, err := c.buildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}

	c.config = cfg

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (c *Centos) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	c.config.Output.Logger = options.Log()

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorUrls, err := c.buildPerVersionMirrorUrls(c.config.Mirrors, c.config.Versions)
	if err != nil {
		return nil, err
	}

	// Execute templated repositories URIs.
	archs := make([]interface{}, 0, len(c.config.Archs))
	for _, v := range c.config.Archs {
		archs = append(archs, interface{}(string(v)))
	}
	err = c.config.BuildTemplates(map[string]interface{}{
		"archs": archs,
	})
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	repositoryURLs, err := c.buildRepositoriesUrls(perVersionMirrorUrls, c.config.Repositories)
	if err != nil {
		return nil, err
	}

	// Get RPM packages from each repository.
	rss := []string{}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}
	searchOptions := rpm.NewSearchOptions(&options, c.config.Archs, rss)
	rpmPackages, err := rpm.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return rpmPackages, nil
}

// Returns the list of version-specific mirror URLs.
func (c *Centos) buildPerVersionMirrorUrls(mirrors []packages.Mirror, versions []distro.Version) ([]*url.URL, error) {
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
func (c *Centos) buildVersions(mirrors []packages.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := c.crawlVersions(mirrors)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (c *Centos) crawlVersions(mirrors []packages.Mirror) ([]distro.Version, error) {
	versions := []distro.Version{}

	seedUrls := make([]*url.URL, 0, len(mirrors))

	for _, mirror := range mirrors {
		u, err := url.Parse(mirror.URL)
		if err != nil {
			return []distro.Version{}, err
		}

		seedUrls = append(seedUrls, u)
	}

	folderNames, err := scrape.CrawlFolders(
		seedUrls,
		CentosMirrorsDistroVersionRegex,
		false,
		c.config.Output.Verbosity >= output.DebugLevel,
	)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

// Returns the list of repositories URLs.
func (c *Centos) buildRepositoriesUrls(roots []*url.URL, repositories []packages.Repository) ([]*url.URL, error) {
	var urls []*url.URL

	for _, root := range roots {
		//nolint:revive,stylecheck
		for _, r := range repositories {
			// Get repository URL from URI.
			//nolint:revive,stylecheck
			us, err := url.JoinPath(root.String(), string(r.URI))
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

// Returns the list of default repositories from the default config.
func (c *Centos) getDefaultRepositories() []packages.Repository {
	var repositories []packages.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
