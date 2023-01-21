package v2

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
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
		Architectures: []packages.Architecture{distro.DefaultArch},
		Versions:      []distro.Version{""},
	}
)
