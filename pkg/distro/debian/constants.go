package debian

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	DebianMirrorsDistroVersionRegex = `^.+$`
)

var (
	// Default configuration for scrape Debian (deb) packages.
	DebianDefaultConfig = distro.Config{
		Mirrors: []packages.Mirror{
			{URL: "https://mirrors.edge.kernel.org/debian/"},
			{URL: "https://security.debian.org/"},
		},
		Repositories: []packages.Repository{
			{Name: "main", URI: packages.URITemplate("/main/")},
			{Name: "contrib", URI: packages.URITemplate("/contrib/")},
			{Name: "non-free", URI: packages.URITemplate("/non-free/")},
		},
		Archs:    []distro.Arch{distro.DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
