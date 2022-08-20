package centos

import (
	"github.com/maxgio92/krawler/pkg/distro"
)

const (
	// Default regex to base the distro version detection on.
	CentosMirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
	CentosPackageFileExtension      = "rpm"
)

var (
	// Default configuration for scrape Centos (RPM) packages.
	//
	// TODO: support for templating default PackagesURI.
	// As of now URI templating depends on distro's viper.AllSettings()
	// data.
	CentosDefaultConfig = distro.Config{
		Mirrors: []distro.Mirror{
			{URL: "https://mirrors.edge.kernel.org/centos/"},
			{URL: "https://archive.kernel.org/centos-vault/"},
		},
		Repositories: []distro.Repository{
			{Name: "base", URI: distro.URITemplate("/os/" + distro.DefaultArch + "/Packages/")},
			{Name: "updates", URI: distro.URITemplate("/updates/" + distro.DefaultArch + "/Packages/")},
			{Name: "BaseOS", URI: distro.URITemplate("/BaseOS/" + distro.DefaultArch + "/os/Packages/")},
			{Name: "AppStream", URI: distro.URITemplate("/AppStream/" + distro.DefaultArch + "/os/Packages/")},
			{Name: "Devel", URI: distro.URITemplate("/Devel/" + distro.DefaultArch + "/os/Packages/")},
		},
		Archs:    []distro.Arch{distro.DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
