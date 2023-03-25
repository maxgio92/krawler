/*
Copyright © 2022 maxgio92 <me@maxgio.it>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fedora

import (
	"net/url"
	"strings"

	"github.com/maxgio92/krawler/pkg/distro"
	"github.com/maxgio92/krawler/pkg/packages"
)

func (f *Fedora) buildConfig(def distro.Config, user distro.Config) (distro.Config, error) {
	config, err := f.mergeConfig(def, user)
	if err != nil {
		return distro.Config{}, err
	}

	err = f.sanitizeConfig(&config)
	if err != nil {
		return distro.Config{}, err
	}

	// Build templated repositories URIs against built-in variables (archs).
	archs := make([]interface{}, 0, len(config.Archs))
	for _, v := range config.Archs {
		archs = append(archs, string(v))
	}
	if err = config.BuildTemplates(map[string]interface{}{
		"archs": archs,
	}); err != nil {
		return distro.Config{}, err
	}

	return config, nil
}

// Returns the final configuration by merging the default with the user provided.
//
//nolint:unparam
func (f *Fedora) mergeConfig(def distro.Config, config distro.Config) (distro.Config, error) {
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
		config.Repositories = f.getDefaultRepositories()
	} else {
		for _, repository := range config.Repositories {
			if repository.URI == "" {
				config.Repositories = f.getDefaultRepositories()

				break
			}
		}
	}

	return config, nil
}

func (f *Fedora) sanitizeConfig(config *distro.Config) error {
	err := f.sanitizeMirrors(&config.Mirrors)
	if err != nil {
		return err
	}

	return nil
}

func (f *Fedora) sanitizeMirrors(mirrors *[]packages.Mirror) error {
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
