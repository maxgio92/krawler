package centos

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
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
		Mirrors: []packages.Mirror{
			{URL: "https://mirrors.edge.kernel.org/centos/"},
			{URL: "https://archive.kernel.org/centos-vault/"},
		},
		Repositories: []packages.Repository{
			{Name: "base", URI: packages.URITemplate("/os/" + distro.DefaultArch + "/")},
			{Name: "updates", URI: packages.URITemplate("/updates/" + distro.DefaultArch + "/")},
			{Name: "BaseOS", URI: packages.URITemplate("/BaseOS/" + distro.DefaultArch + "/os/")},
			{Name: "AppStream", URI: packages.URITemplate("/AppStream/" + distro.DefaultArch + "/os/")},
			{Name: "Devel", URI: packages.URITemplate("/Devel/" + distro.DefaultArch + "/os/")},
		},
		Archs:    []distro.Arch{distro.DefaultArch},
		Versions: nil,
	}

	debugScrape = true
)
