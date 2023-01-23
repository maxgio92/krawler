package amazonlinux

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/scrape"
	"net/url"
)

type AmazonLinux struct {
	Config distro.Config
}

func (a *AmazonLinux) ConfigureCommon(def distro.Config, config distro.Config) error {
	c, err := mergeAndSanitizeConfig(def, config)
	if err != nil {
		return err
	}
	a.Config = c

	return nil
}

// Returns the list of version-specific mirror URLs.
func (a *AmazonLinux) BuildMirrorURLs(mirrors []p.Mirror, versions []distro.Version) ([]*url.URL, error) {
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

// Returns the list of repositories URLs.
func BuildRepositoriesURLs(roots []*url.URL, repositories []p.Repository) ([]*url.URL, error) {
	var urls []*url.URL

	for _, root := range roots {
		for _, r := range repositories {

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

// Returns a list of distro versions, considering the user-provided configuration,
// and if not, the ones available on configured mirrors.
func (a *AmazonLinux) buildVersions(mirrors []p.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := a.crawlVersions(mirrors)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (a *AmazonLinux) crawlVersions(mirrors []p.Mirror) ([]distro.Version, error) {
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
		MirrorsDistroVersionRegex,
		true,
		a.Config.Output.Verbosity >= output.DebugLevel,
	)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}
