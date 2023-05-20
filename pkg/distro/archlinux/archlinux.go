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

package archlinux

import (
	"fmt"
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/alpm"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type ArchLinux struct {
	config distro.Config
}

func (a *ArchLinux) Configure(config distro.Config) error {
	cfg, err := a.buildConfig(DefaultConfig, config)
	if err != nil {
		return err
	}

	a.config = cfg

	return nil
}

// GetPackages scrapes each mirror, for each distro version, for each repository,
// for each architecture, and returns slice of Package and optionally an error.
func (a *ArchLinux) SearchPackages(options packages.SearchOptions) ([]packages.Package, error) {
	a.config.Output.Logger = options.Log()

	// Arch Linux is a rolling release distribution.
	mirrorURLs := []*url.URL{}
	for _, v := range a.config.Mirrors {
		u, err := url.Parse(v.URL)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing mirror URL")
		}

		mirrorURLs = append(mirrorURLs, u)
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	repositoryURLs, err := a.buildRepositoriesUrls(mirrorURLs, a.config.Repositories)
	if err != nil {
		return nil, err
	}

	// Get packages from each repository.
	rss := []string{}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}

	// E.g. https://ftp.halifax.rwth-aachen.de/archlinux/core/os/x86_64/core.db.tar.gz
	res := []packages.Package{}
	for _, v := range rss {
		var repoURL string
		var repo string

		switch {
		case strings.Contains(v, RepoCoreDebug):
			repo = RepoCoreDebug
		case strings.Contains(v, RepoCore):
			repo = RepoCore
		case strings.Contains(v, RepoCommunity):
			repo = RepoCommunity
		case strings.Contains(v, RepoCommunityDebug):
			repo = RepoCommunityDebug
		case strings.Contains(v, RepoCommunityTestingDebug):
			repo = RepoCommunityTestingDebug
		case strings.Contains(v, RepoCommunityTesting):
			repo = RepoCommunityTesting
		case strings.Contains(v, RepoCommunityStagingDebug):
			repo = RepoCommunityStagingDebug
		case strings.Contains(v, RepoCommunityStaging):
			repo = RepoCommunityStaging
		case strings.Contains(v, RepoExtraDebug):
			repo = RepoExtraDebug
		case strings.Contains(v, RepoExtra):
			repo = RepoExtra
		case strings.Contains(v, RepoStagingDebug):
			repo = RepoStagingDebug
		case strings.Contains(v, RepoStaging):
			repo = RepoStaging
		case strings.Contains(v, RepoTestingDebug):
			repo = RepoTestingDebug
		case strings.Contains(v, RepoTesting):
			repo = RepoTesting
		}

		repoURL, err = url.JoinPath(v, fmt.Sprintf("%s.db.tar.gz", repo))
		if err != nil {
			return nil, errors.Wrap(err, "error generating repository URL")
		}

		packageNames := []string{
			options.PackageName(),
			"linux-lts-headers",
			"linux-aarch64-headers",
			"linux-armv7-headers",
		}

		// TODO: remove this serial work.
		// ALPM binding seems to not work with multiple package names.
		for _, v := range packageNames {
			ps, err := alpm.SearchPackages(repoURL, []string{v})
			if err != nil {
				return nil, err
			}

			res = append(res, ps...)
		}
	}

	return res, nil
}

// Returns the list of repositories URLs.
func (a *ArchLinux) buildRepositoriesUrls(roots []*url.URL, repositories []packages.Repository) ([]*url.URL, error) {
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
func (a *ArchLinux) getDefaultRepositories() []packages.Repository {
	var repositories []packages.Repository

	for _, repository := range DefaultConfig.Repositories {
		if !distro.RepositorySliceContains(repositories, repository) {
			repositories = append(repositories, repository)
		}
	}

	return repositories
}
