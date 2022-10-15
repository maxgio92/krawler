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
	"github.com/spf13/cobra"
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
