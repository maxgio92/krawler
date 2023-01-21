package debian

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	DebianMirrorsDistroVersionRegex                       = `^.+$`
	DefaultArch                                           = X86_64Arch
	X86_64Arch                      packages.Architecture = "amd64"
)

var (
	DefaultConfig = distro.Config{
		Mirrors: []packages.Mirror{
			{URL: "https://mirrors.edge.kernel.org/debian/"},
			{URL: "https://security.debian.org/"},
		},
		Repositories: []packages.Repository{
			{Name: "main", URI: packages.URITemplate("/main/")},
			{Name: "contrib", URI: packages.URITemplate("/contrib/")},
			{Name: "non-free", URI: packages.URITemplate("/non-free/")},
		},
		Architectures: []packages.Architecture{DefaultArch},

		// Distribution versions, i.e. Debian dists
		Versions: nil,
	}
)
