package utils

import (
	v "github.com/spf13/viper"

	d "github.com/maxgio92/krawler/pkg/distro"
)

func GetDistroConfigAndVarsFromViper(viper *v.Viper) (d.Config, map[string]interface{}, error) {
	// The distro configuration.
	var config d.Config

	// The distro all settings from Viper
	var allsettings map[string]interface{}

	// The distro config variables from Viper
	var varsSettings map[string]interface{}

	if distros := v.Sub("distros"); distros != nil {
		centos := distros.Sub(string(d.CentosType))
		if centos != nil {
			err := centos.Unmarshal(&config)
			if err != nil {
				return d.Config{}, nil, err
			}

			allsettings = centos.AllSettings()
		}

		amazonLinuxV1 := distros.Sub(string(d.AmazonLinuxV1Type))
		if amazonLinuxV1 != nil {
			err := amazonLinuxV1.Unmarshal(&config)
			if err != nil {
				return d.Config{}, nil, err
			}

			allsettings = amazonLinuxV1.AllSettings()
		}

		amazonLinuxV2 := distros.Sub(string(d.AmazonLinuxV2Type))
		if amazonLinuxV2 != nil {
			err := amazonLinuxV2.Unmarshal(&config)
			if err != nil {
				return d.Config{}, nil, err
			}

			allsettings = amazonLinuxV2.AllSettings()
		}

		amazonLinuxV2022 := distros.Sub(string(d.AmazonLinuxV2022Type))
		if amazonLinuxV2022 != nil {
			err := amazonLinuxV2022.Unmarshal(&config)
			if err != nil {
				return d.Config{}, nil, err
			}

			allsettings = amazonLinuxV2022.AllSettings()
		}
	}

	if _, ok := allsettings["vars"].(map[string]interface{}); ok {
		//nolint:forcetypeassert
		varsSettings = allsettings["vars"].(map[string]interface{})
	}

	vars := MergeMapsAndDeleteKeys(allsettings, varsSettings, "vars", "mirrors", "repositories")

	return config, vars, nil
}
