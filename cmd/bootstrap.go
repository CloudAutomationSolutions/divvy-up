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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bootstrapCmd represents the bootstrap command
var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstrap the provider account with the required resources",
	Long: `This is an easy way to get started.
It requires some elevated priviledges to be executed.
The templates aim to be serverless when possible.
The deployment will take place based on the specified backend.`,
	Run: bootstrap,
}

func init() {
	rootCmd.AddCommand(bootstrapCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// bootstrapCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	bootstrapCmd.Flags().StringP("parameters",
		"p",
		"",
		"These are key value pairs necessary for provisioning the backend provided in yaml format",
	)
	// https://github.com/spf13/viper/issues/397
	// bootstrapCmd.MarkFlagRequired("parameters")
	viper.BindPFlag("parameters", distributeCmd.Flags().Lookup("parameters"))
}

func bootstrap(cmd *cobra.Command, args []string) {
	backendInstance.Bootstrap(viper.GetString("parameters"))
}
