package oracle

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
		{URL: "https://yum.oracle.com/repo/OracleLinux/"},
	},
	Repositories: []packages.Repository{
		{Name: "", URI: packages.URITemplate("/latest/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/MODRHCK/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEK/latest/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR3/latest/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR3/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR4/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR5/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR6/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/UEKR7/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/baseos/latest/{{ .archs }}/")},
		{Name: "", URI: packages.URITemplate("/appstream/{{ .archs }}/")},
	},
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"ppc64le",
	},
	Versions: []distro.Version{
		"OL6",
		"OL7",
		"OL8",
		"OL9",
	},
}
