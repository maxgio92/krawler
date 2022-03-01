/*
Copyright © 2022 maxgio92 <massimiliano.giovagnoli.1992@gmail.com>

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
package cmd

import (
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/maxgio92/krawler/internal/format"
	"github.com/maxgio92/krawler/internal/utils"
	"github.com/maxgio92/krawler/pkg/scraper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// centosCmd represents the centos command.
var centosCmd = &cobra.Command{
	Use:   "centos",
	Short: "List CentOS kernel releases with headers available from mirrors",
	RunE: func(cmd *cobra.Command, args []string) error {
		releases, err := scrape()
		if err != nil {
			return err
		}

		Output, err = format.Encode(Output, releases, format.Type(outputFormat))
		if err != nil {
			return err
		}
		return nil
	},
}

var centosDefaultMirrors = []string{
	"https://mirrors.edge.kernel.org/centos/",
	"https://archive.kernel.org/centos-vault/",
	"https://mirror.nsc.liu.se/centos-store/",
}

func init() {
	listCmd.AddCommand(centosCmd)
}

func scrape() ([]kernelrelease.KernelRelease, error) {
	s, err := scraper.Factory(scraper.Centos)
	if err != nil {
		return nil, err
	}

	// Get packages
	packagePrefix := "kernel-devel"
	mirrorsConfig := scraper.MirrorsConfig{}

	if u := viper.GetStringSlice("mirrors.centos"); u != nil {
		mirrorsConfig.URLs = u
	} else {
		mirrorsConfig.URLs = centosDefaultMirrors
	}

	mirrorsConfig.Archs = []string{"x86_64"}
	mirrorsConfig.PackagesURIFormats = []string{
		"/BaseOS/%s/os/Packages/",
		"/os/%s/Packages/",
		"/updates/%s/Packages/",
	}

	logrus.Debug("Scraping with config: ", mirrorsConfig)

	packages, err := s.Scrape(mirrorsConfig, packagePrefix)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Debug("Getting kernel releases from packages: ", packages)

	// Get kernel releases from kernel header packages
	kernelReleases := []kernelrelease.KernelRelease{}

	for _, v := range packages {
		s := utils.KernelReleaseFromPackageName(v, packagePrefix)
		r := kernelrelease.FromString(s)
		kernelReleases = append(kernelReleases, r)
	}

	logrus.Debug("Obtained KernelRelease objects: ", kernelReleases)

	return utils.UniqueKernelReleases(kernelReleases), nil
}
