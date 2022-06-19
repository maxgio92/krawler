package utils

import (
	"fmt"

	"github.com/maxgio92/krawler/pkg/scrape"
	v "github.com/spf13/viper"
)

//func GetScrapeConfigFromViper(distro string, viper *v.Viper) (scrape.Config, error) {
func GetScrapeConfigFromViper(distro string, viper *v.Viper) (scrape.Config, error) {
	if distro == "" {
		return scrape.Config{}, fmt.Errorf("configuration is not valid")
	}

	var config scrape.Config

	distroViper := viper.Sub(distro)
	distroViper.Unmarshal(&config)

	return config, nil
}
