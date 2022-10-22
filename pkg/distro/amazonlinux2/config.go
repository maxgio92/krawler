package amazonlinux2

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
	"net/url"
	"strings"
)

func (a *AmazonLinux2) buildConfig(def distro.Config, user distro.Config) (distro.Config, error) {
	config, err := a.mergeConfig(def, user)
	if err != nil {
		return distro.Config{}, err
	}

	err = a.sanitizeConfig(&config)
	if err != nil {
		return distro.Config{}, err
	}

	return config, nil
}

// mergeConfig returns the final configuration by merging the default with the user provided.
func (a *AmazonLinux2) mergeConfig(def distro.Config, config distro.Config) (distro.Config, error) {
	if len(config.Archs) < 1 {
		config.Archs = def.Archs
	} else {
		for _, arch := range config.Archs {
			if arch == "" {
				config.Archs = def.Archs

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
		config.Repositories = a.getDefaultRepositories()
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = a.getDefaultRepositories()

				break
			}
		}
	}

	// Force Amazon Linux versions as folder URLs are forbidden.
	if len(config.Versions) < 1 {
		config.Versions = def.Versions
	}

	return config, nil
}

func (a *AmazonLinux2) sanitizeConfig(config *distro.Config) error {
	err := a.sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func (a *AmazonLinux2) sanitizeMirrors(mirrors *[]packages.Mirror) error {
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
