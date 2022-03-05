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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// centosCmd represents the centos command.
var centosCmd = &cobra.Command{
	Use:   "centos",
	Short: "List CentOS kernel releases with headers available from mirrors",
	RunE: func(cmd *cobra.Command, args []string) error {
		kernelReleases, err := getKernelReleases()
		if err != nil {
			return err
		}

		Output, err = format.Encode(Output, kernelReleases, format.Type(outputFormat))
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(centosCmd)
}

func getKernelReleases() ([]kernelrelease.KernelRelease, error) {

	// A representation of a Linux distribution package scraper.
	distro, err := scrape.Factory(scrape.Centos)
	if err != nil {
		return nil, err
	}

	// The filter for filter packages.
	filter := scrape.Filter("kernel-devel")

	// The scraping configuration.
	config, err := utils.getScrapeConfigFromViper(viper)
	if err != nil {
		return err
	}

	// Scrape mirrors for packeges by filter.
	packages, err := distro.GetPackages(config, filter)
	if err != nil {
		return err
	}

	// Get kernel releases from kernel header packages.
	kernelReleases, err := utils.GetKernelReleaseListFromPackageList(packages)
	if err != nil {
		return err
	}

	return kernelReleases
}
