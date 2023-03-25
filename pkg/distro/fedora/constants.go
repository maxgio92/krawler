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
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	// Default regex to base the distro version detection on.
	DistroVersionRegex = `^(0|[1-9]\d*)\/$`
)

var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{Name: "releases", URL: "https://mirrors.edge.kernel.org/fedora/releases/"},
		{Name: "updates", URL: "https://mirrors.edge.kernel.org/fedora/updates/"},
	},
	Repositories: []packages.Repository{
		{Name: "releases", URI: packages.URITemplate("/Everything/{{ .archs }}/os/")},
		{Name: "updates", URI: packages.URITemplate("/Everything/{{ .archs }}/")},
	},
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"armhfp",
		"ppc64le",
	},
	Versions: nil,
}
