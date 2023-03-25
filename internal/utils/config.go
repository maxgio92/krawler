package utils

import (
	v "github.com/spf13/viper"

	d "github.com/maxgio92/krawler/pkg/distro"
)

//nolint:cyclop,funlen,gocognit
func GetDistroConfigAndVarsFromViper(viper *v.Viper) (d.Config, error) {
	// The distro configuration.
	config := d.Config{}

	// The distro all settings from Viper
	var allsettings map[string]interface{}

	// The distro config variables from Viper
	var varsSettings map[string]interface{}

	if output := viper.Sub("output"); output != nil {
		if err := output.Unmarshal(&config.Output); err != nil {
			return d.Config{}, err
		}
	}

	//nolint:nestif
	if distros := viper.Sub("distros"); distros != nil {
		if centos := distros.Sub(d.CentosType); centos != nil {
			if err := centos.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = centos.AllSettings()
		}

		if amazonLinuxV1 := distros.Sub(d.AmazonLinuxV1Type); amazonLinuxV1 != nil {
			if err := amazonLinuxV1.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV1.AllSettings()
		}

		if amazonLinuxV2 := distros.Sub(d.AmazonLinuxV2Type); amazonLinuxV2 != nil {
			if err := amazonLinuxV2.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV2.AllSettings()
		}

		if amazonLinuxV2022 := distros.Sub(d.AmazonLinuxV2022Type); amazonLinuxV2022 != nil {
			if err := amazonLinuxV2022.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV2022.AllSettings()
		}

		if debian := distros.Sub(d.DebianType); debian != nil {
			if err := debian.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = debian.AllSettings()
		}

		if ubuntu := distros.Sub(d.UbuntuType); ubuntu != nil {
			if err := ubuntu.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = ubuntu.AllSettings()
		}

		if fedora := distros.Sub(d.FedoraType); fedora != nil {
			if err := fedora.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = fedora.AllSettings()
		}
	}

	if _, ok := allsettings["vars"].(map[string]interface{}); ok {
		//nolint:forcetypeassert
		varsSettings = allsettings["vars"].(map[string]interface{})
	}

	vars := MergeMapsAndDeleteKeys(allsettings, varsSettings, "vars", "mirrors", "repositories")

	err := config.BuildTemplates(vars)
	if err != nil {
		return d.Config{}, err
	}

	return config, nil
}
