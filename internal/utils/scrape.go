package utils

import (
	"fmt"

	"github.com/maxgio92/krawler/pkg/distro"
	v "github.com/spf13/viper"
)

//func GetScrapeConfigFromViper(distro string, viper *v.Viper) (scrape.Config, error) {
func GetScrapeConfigFromViper(distroName string, viper *v.Viper) (distro.Config, error) {
	if distroName == "" {
		return distro.Config{}, fmt.Errorf("configuration is not valid")
	}

	var config distro.Config

	distroViper := viper.Sub(distroName)
	distroViper.Unmarshal(&config)

	return config, nil
}
