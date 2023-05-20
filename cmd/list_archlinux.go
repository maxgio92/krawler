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
	"github.com/maxgio92/krawler/internal/format"
	"github.com/maxgio92/krawler/pkg/distro/archlinux"
	"github.com/spf13/cobra"
)

// fedoraCmd represents the fedora command.
var archLinuxCmd = &cobra.Command{
	Use:   "archlinux",
	Short: "List Arch Linux kernel releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		kernelReleases, err := getKernelReleases(&archlinux.ArchLinux{}, "linux-headers")
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
	listCmd.AddCommand(archLinuxCmd)
}
