// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/cloudautomationsolutions/divvy-up/pkg/provider"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var backendInstance provider.Backend

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "divvy-up",
	Short: "a tool that is meant for sharing secret values using your own cloud provider of choice, as a backend",
	Long: `This project aims to be a safe and secure way of sharing secrets between individuals in a team.
	
The goal is to stop people sending secrets over slack.
It is designed to be included in your own infrstructure be it AWS, Google Compute Cloud, Kubernetes with support for other platforms to come.
The sensitive data that needs to be shared will be stored in the backend for temporary use.
The aim is to enable as much auditing as possible. The file should be kept temporary which means the ideal way of handling a download is to delete it afterwards.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		backendInstance = provider.BackendFromFlag(viper.GetString("backend"))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// This flag points to the configuration file of the application
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.divvy-up.yaml)")

	// This flag chooses the backend to be used
	rootCmd.PersistentFlags().StringP("backend", "b", "", "backend to be used")
	rootCmd.MarkFlagRequired("backend")

	viper.BindPFlag("backend", rootCmd.PersistentFlags().Lookup("backend"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".divvy-up" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".divvy-up")
	}
	// viper.SetEnvPrefix

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
