package amazonlinux2

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

const (
	MirrorsDistroVersionRegex = `^(0|[1-9]\d*)(\.(0|[1-9]\d*)?)?(\.(0|[1-9]\d*)?)?(-[a-zA-Z\d][-a-zA-Z.\d]*)?(\+[a-zA-Z\d][-a-zA-Z.\d]*)?\/$`
	PackageFileExtension      = "rpm"
)

var (
	// DefaultConfig is the default configuration for scrape Amazon Linux (RPM) packages.
	// As of now URI templating depends on distro's viper.AllSettings() data.
	DefaultConfig = distro.Config{
		Mirrors: []packages.Mirror{
			{
				Name: "AL2",
				URL:  "http://amazonlinux.us-east-1.amazonaws.com/2/",
			},
		},
		Repositories: []packages.Repository{
			{Name: "", URI: "core/2.0"},
			{Name: "", URI: "core/latest"},
			{Name: "", URI: "extras/kernel-5.4/latest"},
			{Name: "", URI: "extras/kernel-5.10/latest"},
			{Name: "", URI: "extras/kernel-5.15/latest"},
		},
		Archs:    []distro.Arch{distro.DefaultArch},
		Versions: []distro.Version{""},
	}

	debugScrape = true
)
