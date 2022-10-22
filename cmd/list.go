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
	"github.com/maxgio92/krawler/internal/utils"
	"github.com/maxgio92/krawler/pkg/distro"
	kr "github.com/maxgio92/krawler/pkg/kernelrelease"
	"github.com/maxgio92/krawler/pkg/packages"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

var (
	// The output format flag value.
	outputFormat string

	// listCmd represents the list command.
	listCmd = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available kernel releases with distributed headers, by Linux distribution",
	}
)

func init() {
	rootCmd.AddCommand(listCmd)

	// Bind the output format flag. Default is text.
	listCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json, yaml)")
}

func getKernelReleases(distro distro.Distro) ([]kr.KernelRelease, error) {
	// The filter for filter packages.
	filter := packages.NewFilter(KernelHeadersPackageName, ".config")

	config, vars, err := utils.GetDistroConfigAndVarsFromViper(v.GetViper())
	if err != nil {
		return []kr.KernelRelease{}, err
	}

	err = distro.Configure(config, vars)
	if err != nil {
		return []kr.KernelRelease{}, err
	}

	// Scrape mirrors for packeges by filter.
	packages, err := distro.GetPackages(*filter)
	if err != nil {
		return []kr.KernelRelease{}, err
	}

	// Get kernel releases from kernel header packages.
	kernelReleases, err := kr.GetKernelReleaseListFromPackageList(packages, KernelHeadersPackageName)
	if err != nil {
		return []kr.KernelRelease{}, err
	}

	return kernelReleases, nil
}
