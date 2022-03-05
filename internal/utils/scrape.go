package utils

import (
	"github.com/maxgio92/krawler/pkg/scrape"
	"github.com/spf13/viper"
)

func getScrapeConfigFromViper(distro string, viper viper.Viper) (scrape.Config, error) {
	config := scrape.Config{}

	if distro != nil {
		if mirrors := viper.GetStringSlice("distros."+distro+".mirrors"); mirrors != nil {
			for mirror in range mirrors {
				if mirror["baseUrl"] != nil {
					config.Mirrors = append(config.Mirrors, scrape.Mirror.BaseUrl(mirror["baseUrl"]))
				}
			}
		}

		if archs := viper.GetStringSlice("distros."+distro+".archs"); archs != nil {
			for arch in range archs {
				if arch["id"] != nil {
					config.Archs = append(config.Archs, scrape.Arch(arch["id"]))
				}
			}
		}
	}

	return config, nil
}
