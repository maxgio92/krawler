package debian

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	DebianMirrorsDistroVersionRegex                       = `^.+$`
	DefaultArch                                           = X8664Arch
	X8664Arch                       packages.Architecture = "amd64"
)

var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{URL: "https://mirrors.edge.kernel.org/debian/"},
		{URL: "http://security.debian.org"},
	},
	Repositories: []packages.Repository{
		{Name: "main", URI: packages.URITemplate("main")},
		{Name: "contrib", URI: packages.URITemplate("contrib")},
		{Name: "non-free", URI: packages.URITemplate("non-free")},
		{Name: "multiverse", URI: packages.URITemplate("multiverse")},
		{Name: "universe", URI: packages.URITemplate("universe")},
		{Name: "restricted", URI: packages.URITemplate("restricted")},
	},
	Archs: []packages.Architecture{DefaultArch},

	// Distribution versions, i.e. Debian dists
	Versions: nil,
}
