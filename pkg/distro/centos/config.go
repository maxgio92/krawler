package centos

import (
	"net/url"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

func (c *Centos) buildConfig(def distro.Config, user distro.Config) (distro.Config, error) {
	config, err := c.mergeConfig(def, user)
	if err != nil {
		return distro.Config{}, err
	}

	err = c.sanitizeConfig(&config)
	if err != nil {
		return distro.Config{}, err
	}

	return config, nil
}

// Returns the final configuration by merging the default with the user provided.
//
//nolint:unparam
func (c *Centos) mergeConfig(def distro.Config, config distro.Config) (distro.Config, error) {
	if len(config.Architectures) < 1 {
		config.Architectures = def.Architectures
	} else {
		for _, arch := range config.Architectures {
			if arch == "" {
				config.Architectures = def.Architectures

				break
			}
		}
	}

	if len(config.Mirrors) < 1 {
		config.Mirrors = def.Mirrors
	} else {
		for _, mirror := range config.Mirrors {
			if mirror.URL == "" {
				config.Mirrors = def.Mirrors

				break
			}
		}
	}

	if len(config.Repositories) < 1 {
		config.Repositories = c.getDefaultRepositories()
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = c.getDefaultRepositories()

				break
			}
		}
	}

	return config, nil
}

func (c *Centos) sanitizeConfig(config *distro.Config) error {
	err := c.sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func (c *Centos) sanitizeMirrors(mirrors *[]packages.Mirror) error {
	for i, mirror := range *mirrors {
		if !strings.HasSuffix(mirror.URL, "/") {
			(*mirrors)[i].URL = mirror.URL + "/"
		}

		_, err := url.Parse(mirror.URL)
		if err != nil {
			return err
		}
	}

	return nil
}
