package centos

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	// Default regex to base the distro version detection on.
	CentosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
)

var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{URL: "https://mirrors.edge.kernel.org/centos/"},
		{URL: "https://archive.kernel.org/centos-vault/"},
	},
	Repositories: []packages.Repository{
		{Name: "base", URI: packages.URITemplate("/os/{{ .archs }}/")},
		{Name: "updates", URI: packages.URITemplate("/updates/{{ .archs }}/")},
		{Name: "BaseOS", URI: packages.URITemplate("/BaseOS/{{ .archs }}/os/")},
		{Name: "AppStream", URI: packages.URITemplate("/AppStream/{{ .archs }}/os/")},
		{Name: "Devel", URI: packages.URITemplate("/Devel/{{ .archs }}/os/")},
	},
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"ppc64le",
	},
	Versions: nil,
}
