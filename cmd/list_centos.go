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
package cmd

import (
	"github.com/falcosecurity/driverkit/pkg/kernelrelease"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"

	"github.com/maxgio92/krawler/internal/format"
	"github.com/maxgio92/krawler/internal/utils"
	d "github.com/maxgio92/krawler/pkg/distro"
)

// centosCmd represents the centos command.
var centosCmd = &cobra.Command{
	Use:   "centos",
	Short: "List CentOS kernel releases with headers available from mirrors",
	RunE: func(cmd *cobra.Command, args []string) error {
		kernelReleases, err := getKernelReleases()
		cobra.CheckErr(err)

		if len(kernelReleases) > 0 {
			Output, err = format.Encode(Output, kernelReleases, format.Type(outputFormat))
			cobra.CheckErr(err)
		} else {
			//nolint:errcheck
			Output.WriteString("No releases found.\n")
		}

		return nil
	},
}

func init() {
	listCmd.AddCommand(centosCmd)
}

func getKernelReleases() ([]kernelrelease.KernelRelease, error) {
	// A representation of a Linux distribution package scraper.
	distro, err := d.NewDistro(d.CentosType)
	if err != nil {
		return nil, err
	}

	// The filter for filter packages.
	packagePrefix := KernelHeadersPackageName
	filter := d.Filter(packagePrefix)

	config, vars, err := utils.GetDistroConfigAndVarsFromViper(v.GetViper())
	if err != nil {
		return []kernelrelease.KernelRelease{}, err
	}

	err = distro.Configure(config, vars)
	if err != nil {
		return []kernelrelease.KernelRelease{}, err
	}

	// Scrape mirrors for packeges by filter.
	packages, err := distro.GetPackages(filter)
	if err != nil {
		return []kernelrelease.KernelRelease{}, err
	}

	// Get kernel releases from kernel header packages.
	kernelReleases, err := utils.GetKernelReleaseListFromPackageList(packages, packagePrefix)
	if err != nil {
		return []kernelrelease.KernelRelease{}, err
	}

	return kernelReleases, nil
}
