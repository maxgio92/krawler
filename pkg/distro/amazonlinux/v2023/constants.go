package v2023

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

// DefaultConfig is the default configuration for scrape Amazon Linux (RPM) packages.
// As of now URI templating depends on distro's viper.AllSettings() data.
var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{
			Name: "AL2023",
			URL:  "https://cdn.amazonlinux.com/al2023/core/mirrors/",
		},
	},
	Repositories: []packages.Repository{
		{Name: "", URI: "latest"},
	},
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"ppc64le",
	},
	Versions: []distro.Version{""},
}
