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
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	RepoCore                  = "core"
	RepoCoreDebug             = "core-debug"
	RepoCommunity             = "community"
	RepoCommunityTesting      = "community-testing"
	RepoCommunityStaging      = "community-staging"
	RepoCommunityDebug        = "community-debug"
	RepoCommunityTestingDebug = "community-testing-debug"
	RepoCommunityStagingDebug = "community-staging-debug"
	RepoExtra                 = "extra"
	RepoExtraDebug            = "extra-debug"
	RepoStaging               = "staging"
	RepoStagingDebug          = "staging-debug"
	RepoTesting               = "testing"
	RepoTestingDebug          = "testing-debug"

	archiveMonthRetention    = 3
	archiveReleaseDayOfMonth = 10
	archiveMirror            = "https://archive.archlinux.org/repos/"
)

var (
	// As Arch Linux is a rolling release distribution, we need archives to track previous rollouts.
	repos = []string{
		RepoCore,
		RepoCoreDebug,
		RepoCommunity,
		RepoCommunityTesting,
		RepoCommunityDebug,
		RepoCommunityTestingDebug,
		RepoCommunityStagingDebug,
		RepoExtra,
		RepoExtraDebug,
		RepoStaging,
		RepoStagingDebug,
		RepoTesting,
		RepoTestingDebug,
	}

	DefaultConfig = distro.Config{
		Mirrors: []packages.Mirror{
			{Name: "arm64", URL: "http://de.mirror.archlinuxarm.org/aarch64/"},
			{Name: "arm64", URL: "http://de.mirror.archlinuxarm.org/aarch64/"},
			{Name: "arm32", URL: "http://de.mirror.archlinuxarm.org/armv7h/"},
			{Name: "archmirror", URL: "https://archmirror.it/repos/"},
			{Name: "kernel.org", URL: "https://mirrors.edge.kernel.org/archlinux/"},
		},
		Repositories: []packages.Repository{
			{Name: "core", URI: packages.URITemplate("/core/os/{{ .archs }}/")},
			{Name: "aur", URI: packages.URITemplate("/aur/os/{{ .archs }}/")},
			{Name: "community", URI: packages.URITemplate("/community/os/{{ .archs }}/")},
			{Name: "extra", URI: packages.URITemplate("/extra/os/{{ .archs }}/")},

			// Architecture is embedded already in the mirror URL.
			{Name: "core-arm", URI: packages.URITemplate("/core-arm/")},
			{Name: "aur-arm", URI: packages.URITemplate("/aur-arm/")},
			{Name: "community-arm", URI: packages.URITemplate("/community-arm/")},
			{Name: "extra-arm", URI: packages.URITemplate("/extra-arm/")},
		},
		Archs: []packages.Architecture{
			"x86_64",
			"aarch64",
			"armv7h",
		},

		// Arch Linux is a rollin-release distribution.
		Versions: nil,
	}

	archiveMirrorURLs = []string{archiveMirror}
	archiveRepos      = repos
)
