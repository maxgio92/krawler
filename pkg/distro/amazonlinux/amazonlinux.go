package amazonlinux

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	p "github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"github.com/maxgio92/krawler/pkg/scrape"
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

// BuildMirrorURLs returns the list of version-specific mirror URLs.
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

// BuildRepositoryURLs returns the list of repositories URLs.
func BuildRepositoryURLs(roots []*url.URL, repositories []p.Repository) ([]*url.URL, error) {
	var urls []*url.URL

	for _, root := range roots {
		for _, r := range repositories {
			us, err := url.JoinPath(root.String(), string(r.URI))
			if err != nil {
				return nil, err
			}

			repositoryURL, err := url.Parse(us)
			if err != nil {
				return nil, err
			}

			urls = append(urls, repositoryURL)
		}
	}

	return urls, nil
}

// buildVersions returns a list of distro versions, considering the user-provided configuration,
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

// crawlVersions returns the list of the current available distro versions, by scraping
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

// SearchPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (a *AmazonLinux) SearchPackages(options p.SearchOptions) ([]p.Package, error) {
	a.Config.Output.Logger = options.Log()

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorURLs, err := a.BuildMirrorURLs(a.Config.Mirrors, a.Config.Versions)
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	repositoriesURLrefs, err := BuildRepositoryURLs(perVersionMirrorURLs, a.Config.Repositories)
	if err != nil {
		return nil, err
	}

	// Dereference repository URLs.
	repositoryURLs, err := a.dereferenceRepositoryURLs(repositoriesURLrefs, a.Config.Archs)
	if err != nil {
		return nil, err
	}

	// Get RPM packages from each repository.
	rss := []string{}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}

	searchOptions := rpm.NewSearchOptions(&options, a.Config.Archs, rss)
	rpmPackages, err := rpm.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return rpmPackages, nil
}

func (a *AmazonLinux) dereferenceRepositoryURLs(repoURLs []*url.URL, archs []p.Architecture) ([]*url.URL, error) {
	var urls []*url.URL

	for _, ar := range archs {
		for _, v := range repoURLs {
			r, err := a.dereferenceRepositoryURL(v, ar)
			if err != nil {
				return nil, err
			}

			if r != nil {
				urls = append(urls, r)
			}
		}
	}

	return urls, nil
}

func (a *AmazonLinux) dereferenceRepositoryURL(src *url.URL, arch p.Architecture) (*url.URL, error) {
	var dest *url.URL

	mirrorListURL, err := url.JoinPath(src.String(), string(arch), "mirror.list")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, mirrorListURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		a.Config.Output.Logger.Error("Amazon Linux v2023 repository URL not valid to be dereferenced")
		//nolint:nilnil
		return nil, nil
	}

	if resp.Body == nil {
		a.Config.Output.Logger.Error("empty response from Amazon Linux v2023 repository reference URL")
		//nolint:nilnil
		return nil, nil
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Get first repository URL available, no matter what the geolocation.
	s := strings.Split(string(b), "\n")[0]

	dest, err = url.Parse(s)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
