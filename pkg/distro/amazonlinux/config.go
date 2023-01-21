package amazonlinux

import (
	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
	"net/url"
	"strings"
)

func mergeAndSanitizeConfig(def distro.Config, user distro.Config) (distro.Config, error) {
	config, err := mergeConfig(def, user)
	if err != nil {
		return distro.Config{}, err
	}

	err = sanitizeConfig(&config)
	if err != nil {
		return distro.Config{}, err
	}

	return config, nil
}

// mergeConfig returns the final configuration by merging the default with the user provided.
func mergeConfig(def distro.Config, config distro.Config) (distro.Config, error) {
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
		config.Repositories = def.Repositories
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = def.Repositories

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

func sanitizeConfig(config *distro.Config) error {
	err := sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func sanitizeMirrors(mirrors *[]packages.Mirror) error {
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
