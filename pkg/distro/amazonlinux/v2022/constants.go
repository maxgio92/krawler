package v2022

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

// DefaultConfig is the default configuration for scrape Amazon Linux (RPM) packages.
// As of now URI templating depends on distro's viper.AllSettings() data.
var DefaultConfig = distro.Config{
	Mirrors: []packages.Mirror{
		{
			Name: "AL2022",
			URL:  "https://al2022-repos-us-east-1-9761ab97.s3.dualstack.us-east-1.amazonaws.com/core/mirrors/",
		},
	},
	Repositories: []packages.Repository{
		{Name: "", URI: "2022.0.20220202"},
		{Name: "", URI: "2022.0.20220315"},
		{Name: "", URI: "2022.0.20221012"},
	},
	Archs: []packages.Architecture{
		"aarch64",
		"x86_64",
		"ppc64le",
	},
	Versions: []distro.Version{""},
}
