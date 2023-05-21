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
	"net/url"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/packages/alpm"
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/pkg/errors"
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

	mirrorURLs, err := a.buildMirrorURLs()
	if err != nil {
		return nil, errors.Wrap(err, "error building mirror URLs")
	}

	// Build available repository URLs based on provided configuration,
	// for each distribution version.
	rss := []string{}

	repositoryURLs, err := a.buildRepositoriesURLs(mirrorURLs, a.config.Repositories)
	if err != nil {
		return nil, err
	}
	for _, ru := range repositoryURLs {
		rss = append(rss, ru.String())
	}

	// TODO: input all repository URLs.
	archiveURLs, err := a.buildArchiveURLs(archiveMirrorURLs, archiveRepos)
	for _, au := range archiveURLs {
		rss = append(rss, au.String())
	}

	// Get packages from each repository.
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

func (a *ArchLinux) buildMirrorURLs() ([]*url.URL, error) {

	// Arch Linux is a rolling release distribution.
	mirrorURLs := []*url.URL{}
	for _, v := range a.config.Mirrors {
		u, err := url.Parse(v.URL)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing mirror URL")
		}

		mirrorURLs = append(mirrorURLs, u)
	}

	return mirrorURLs, nil
}

// Returns the list of repositories URLs.
func (a *ArchLinux) buildRepositoriesURLs(roots []*url.URL, repositories []packages.Repository) ([]*url.URL, error) {
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

func (a *ArchLinux) buildArchiveURLs(archiveMirrorURLs []string, repositoryNames []string) ([]*url.URL, error) {
	seeds := []*url.URL{}
	for _, v := range archiveMirrorURLs {
		s, err := url.Parse(v)
		if err != nil {
			return nil, err
		}

		seeds = append(seeds, s)
	}

	au, err := scrape.CrawlFoldersPath(seeds, `^core\/$`, true, true)
	if err != nil {
		return nil, errors.Wrap(err, "error building archive repository URLs")
	}

	archiveURLs := []*url.URL{}
	for _, v := range au {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}

		archiveURLs = append(archiveURLs, u)
	}

	return archiveURLs, nil
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
