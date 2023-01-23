package utils

import (
	d "github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/maxgio92/krawler/pkg/utils/template"

	v "github.com/spf13/viper"
)

//nolint:cyclop,funlen,gocognit
func GetDistroConfigAndVarsFromViper(viper *v.Viper) (d.Config, error) {
	// The distro configuration.
	var config d.Config

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
		centos := distros.Sub(d.CentosType)
		if centos != nil {
			if err := centos.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = centos.AllSettings()
		}

		amazonLinuxV1 := distros.Sub(d.AmazonLinuxV1Type)
		if amazonLinuxV1 != nil {
			if err := amazonLinuxV1.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV1.AllSettings()
		}

		amazonLinuxV2 := distros.Sub(d.AmazonLinuxV2Type)
		if amazonLinuxV2 != nil {
			if err := amazonLinuxV2.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV2.AllSettings()
		}

		amazonLinuxV2022 := distros.Sub(d.AmazonLinuxV2022Type)
		if amazonLinuxV2022 != nil {
			if err := amazonLinuxV2022.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = amazonLinuxV2022.AllSettings()
		}

		debian := distros.Sub(d.DebianType)
		if debian != nil {
			if err := debian.Unmarshal(&config); err != nil {
				return d.Config{}, err
			}

			allsettings = debian.AllSettings()
		}
	}

	if _, ok := allsettings["vars"].(map[string]interface{}); ok {
		//nolint:forcetypeassert
		varsSettings = allsettings["vars"].(map[string]interface{})
	}

	vars := MergeMapsAndDeleteKeys(allsettings, varsSettings, "vars", "mirrors", "repositories")

	err := expandTemplatedDistroConfig(&config, vars)
	if err != nil {
		return d.Config{}, err
	}

	return config, nil
}

func expandTemplatedDistroConfig(config *d.Config, vars map[string]interface{}) error {
	repositories := config.Repositories

	uris := []string{}

	for _, repository := range repositories {
		if repository.URI != "" {
			result, err := template.MultiplexAndExecute(string(repository.URI), vars)
			if err != nil {
				return err
			}

			uris = append(uris, result...)
		}
	}

	r := []packages.Repository{}
	for _, v := range uris {
		r = append(r, packages.Repository{Name: "", URI: packages.URITemplate(v)})
	}

	config.Repositories = r

	return nil
}
