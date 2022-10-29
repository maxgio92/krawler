package v1

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

var (
	// DefaultConfig is the default configuration for scrape Amazon Linux (RPM) packages.
	// As of now URI templating depends on distro's viper.AllSettings() data.
	DefaultConfig = distro.Config{
		Mirrors: []packages.Mirror{
			{Name: "AL1", URL: "http://repo.us-east-1.amazonaws.com/"},
		},
		Repositories: []packages.Repository{
			{Name: "", URI: "/updates/"},
			{Name: "", URI: "/main/"},
		},
		Archs: []distro.Arch{distro.DefaultArch},
		Versions: []distro.Version{
			"latest",
			"2017.03",
			"2017.09",
			"2018.03",
		},
	}
)
