package v2

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

// DefaultConfig is the default configuration for scrape Amazon Linux (RPM) packages.
// As of now URI templating depends on distro's viper.AllSettings() data.
var DefaultConfig = distro.Config{
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
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"ppc64le",
	},
	Versions: []distro.Version{""},
}
