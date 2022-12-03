package amazonlinux

import (
	"github.com/maxgio92/krawler/pkg/distro"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/maxgio92/krawler/pkg/utils/template"
	"net/url"
)

type AmazonLinux struct {
	Config distro.Config
	Vars   map[string]interface{}
}

func (a *AmazonLinux) ConfigureCommon(config distro.Config, vars map[string]interface{}) error {
	a.Config = config
	a.Vars = vars

	return nil
}

// Returns the list of version-specific mirror URLs.
func BuildMirrorURLs(mirrors []p.Mirror, versions []distro.Version) ([]*url.URL, error) {
	versions, err := buildVersions(mirrors, versions)
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
func BuildRepositoriesURLs(roots []*url.URL, repositories []p.Repository, vars map[string]interface{}) ([]*url.URL, error) {
	var urls []*url.URL

	uris, err := buildRepositoriesURIs(repositories, vars)
	if err != nil {
		return []*url.URL{}, err
	}

	for _, root := range roots {
		for _, uri := range uris {

			us, err := url.JoinPath(root.String(), uri)
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
func buildVersions(mirrors []p.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := crawlVersions(mirrors, debugScrape)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func crawlVersions(mirrors []p.Mirror, debug bool) ([]distro.Version, error) {
	versions := []distro.Version{}

	seedUrls := make([]*url.URL, 0, len(mirrors))

	for _, mirror := range mirrors {
		u, err := url.Parse(mirror.URL)
		if err != nil {
			return []distro.Version{}, err
		}

		seedUrls = append(seedUrls, u)
	}

	folderNames, err := scrape.CrawlFolders(seedUrls, MirrorsDistroVersionRegex, true, debug)
	if err != nil {
		return []distro.Version{}, err
	}

	for _, v := range folderNames {
		versions = append(versions, distro.Version(v))
	}

	return versions, nil
}

func buildRepositoriesURIs(repositories []p.Repository, vars map[string]interface{}) ([]string, error) {
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
	if len(uris) == 0 {
		uris = append(uris, "/")
	}

	return uris, nil
}
