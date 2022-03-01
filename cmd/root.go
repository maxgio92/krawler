/*
Copyright © 2022 maxgio92 massimiliano.giovagnoli.1992@gmail.com

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
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// The config file flag value.
	cfgFile string

	// The commands output buffer.
	Output = bufio.NewWriter(os.Stdout)

	// The verbose flag value.
	verbosity string

	// rootCmd represents the base command when called without any subcommands.
	rootCmd = &cobra.Command{
		Use:   "krawler",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if err := Output.Flush(); err != nil {
				return fmt.Errorf("cannot flush output: %w", err)
			}

			return nil
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here is where we define the PreRun func, using the verbose flag value.
	// We use the standard output for logs.
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := initLogs(os.Stdout, verbosity); err != nil {
			return err
		}

		return nil
	}

	// Bind the config file flag. Default value is $HOME/.krawler.yaml.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.krawler.yaml)")

	// Bind the verbose flag. Default value is the warn level.
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", logrus.WarnLevel.String(), "Log level (debug, info, warn, error, fatal, panic)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".krawler" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".krawler")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
	}
}

// setUpLogs set the log output ans the log level.
func initLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logrus.SetLevel(lvl)

	return nil
}
