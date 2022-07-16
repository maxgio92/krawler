package utils

import (
	d "github.com/maxgio92/krawler/pkg/distro"

	v "github.com/spf13/viper"
)

func GetDistroConfigAndVarsFromViper(viper *v.Viper) (d.Config, map[string]interface{}, error) {

	// The distro configuration.
	var config d.Config

	// The distro all settings from Viper
	var allsettings map[string]interface{}
	var varsSettings map[string]interface{}

	distros := v.Sub("distros")
	if distros != nil {
		centos := distros.Sub(string(d.CentosType))
		if centos != nil {
			err := centos.Unmarshal(&config)
			if err != nil {
				return d.Config{}, nil, err
			}

			allsettings = centos.AllSettings()
		}
	}

	if _, ok := allsettings["vars"].(map[string]interface{}); ok {
		varsSettings = allsettings["vars"].(map[string]interface{})
	}

	vars := MergeMapsAndDeleteKeys(allsettings, varsSettings, "vars", "mirrors", "repositories")

	return config, vars, nil
}
