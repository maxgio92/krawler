/*
Copyright © 2022 maxgio92 <me@maxgio.it>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fedora

import (
	"net/url"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/output"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/rpm"
	"github.com/maxgio92/krawler/pkg/scrape"
)

type Fedora struct {
	config distro.Config
}

func (f *Fedora) Configure(config distro.Config) error {
	cfg, err := f.buildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}

	f.config = cfg

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (f *Fedora) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	f.config.Output.Logger = options.Log()

	// Build distribution version-specific mirror root URLs.
	perVersionMirrorUrls, err := f.buildPerVersionMirrorUrls(f.config.Mirrors, f.config.Versions)
	if err != nil {
		return nil, err
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	repositoryURLs, err := f.buildRepositoriesUrls(perVersionMirrorUrls, f.config.Repositories)
	if err != nil {
		return nil, err
	}

	// Get RPM packages from each repository.
	rss := []string{}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}
	searchOptions := rpm.NewSearchOptions(&options, f.config.Archs, rss)
	rpmPackages, err := rpm.SearchPackages(searchOptions)
	if err != nil {
		return nil, err
	}

	return rpmPackages, nil
}

// Returns the list of version-specific mirror URLs.
func (f *Fedora) buildPerVersionMirrorUrls(mirrors []packages.Mirror, versions []distro.Version) ([]*url.URL, error) {
	versions, err := f.buildVersions(mirrors, versions)
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
func (f *Fedora) buildVersions(mirrors []packages.Mirror, staticVersions []distro.Version) ([]distro.Version, error) {
	if staticVersions != nil {
		return staticVersions, nil
	}

	var dynamicVersions []distro.Version

	dynamicVersions, err := f.crawlVersions(mirrors)
	if err != nil {
		return nil, err
	}

	return dynamicVersions, nil
}

// Returns the list of the current available distro versions, by scraping
// the specified mirrors, dynamically.
func (f *Fedora) crawlVersions(mirrors []packages.Mirror) ([]distro.Version, error) {
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
		DistroVersionRegex,
		false,
		f.config.Output.Verbosity >= output.DebugLevel,
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
func (f *Fedora) buildRepositoriesUrls(roots []*url.URL, repositories []packages.Repository) ([]*url.URL, error) {
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
func (f *Fedora) getDefaultRepositories() []packages.Repository {
	var repositories []packages.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
