/*
Copyright © 2022 maxgio92 <me@maxgio.it>

Licensed under the Apache License, Version v2.0 (the "License");
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
	"github.com/spf13/cobra"

	"github.com/maxgio92/krawler/internal/format"
	v2023 "github.com/maxgio92/krawler/pkg/distro/amazonlinux/v2023"
)

// amazonLinux2Cmd represents the centos command.
var amazonLinux2023Cmd = &cobra.Command{
	Use:   "amazonlinux2023",
	Short: "List Amazon Linux 2023 kernel releases",
	RunE: func(cmd *cobra.Command, args []string) error {
		kernelReleases, err := getKernelReleases(&v2023.AmazonLinux{}, RPMKernelHeadersPackageName)
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
	listCmd.AddCommand(amazonLinux2023Cmd)
}
