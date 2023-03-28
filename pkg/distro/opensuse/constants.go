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

package opensuse

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	// Default regex to base the distro version detection on.
	// Match both SemVer and tags like 'openSUSE-stable'.
	DistroVersionRegex = `^.+\/$`
)

var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{Name: "default", URL: "https://mirrors.edge.kernel.org/opensuse/distribution/"},
		{Name: "tumbleweed", URL: "https://mirrors.edge.kernel.org/opensuse/"},
		{Name: "leap", URL: "https://mirrors.edge.kernel.org/opensuse/distribution/leap/"},
		{Name: "kernel", URL: "http://download.opensuse.org/repositories/Kernel:/"},
	},
	Repositories: []packages.Repository{
		{Name: "default", URI: packages.URITemplate("/repo/oss/")},
		{Name: "kernel-arm", URI: packages.URITemplate("/ARM/")},
		{Name: "kernel-ppc", URI: packages.URITemplate("/PPC/")},
		{Name: "kernel-riscv", URI: packages.URITemplate("/RISCV/")},
		{Name: "kernel-s390", URI: packages.URITemplate("/S390/")},
		{Name: "kernel-standard", URI: packages.URITemplate("/standard/")},
		{Name: "kernel-ports", URI: packages.URITemplate("/ports/")},
		{Name: "kernel-backport-standard", URI: packages.URITemplate("/Backport/standard")},
		{Name: "kernel-backport-ports", URI: packages.URITemplate("/Backport/ports")},
		{Name: "kernel-submit-standard", URI: packages.URITemplate("/Submit/standard/")},
		{Name: "kernel-submit-ports", URI: packages.URITemplate("/Submit/ports/")},
	},
	Archs: []packages.Architecture{
		"armv6hl",
		"armv7hl",
		"aarch64",
		"armhfp",
		"x86_64",
		"noarch",
		"i686",
		"ppc",
		"ppc64",
		"ppc64le",
		"s390x",
	},

	// Crawl all versions by default, filtering names on the DistroVersionRegex regular expression.
	Versions: nil,
}
