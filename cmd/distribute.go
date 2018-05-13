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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// distributeCmd represents the distribute command
var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "Share a file in a secure way using your cloud provider",
	Long: `This is the main function for enabling sharing.
It points to a file containing your sensitive data and sends it via a secure way to the backend for temporary storage`,
	Run: distribute,
}

func init() {
	rootCmd.AddCommand(distributeCmd)

	distributeCmd.Flags().StringP("file",
		"f",
		"",
		"The file which holds your secrets",
	)
	distributeCmd.Flags().IntP("expiration",
		"e",
		1800,
		"The file which holds your secrets",
	)
	// https://github.com/spf13/viper/issues/397
	// distributeCmd.MarkFlagRequired("file")

	viper.BindPFlag("file", distributeCmd.Flags().Lookup("file"))
	viper.BindPFlag("expiration", distributeCmd.Flags().Lookup("expiration"))
}

func distribute(cmd *cobra.Command, args []string) {
	if !viper.IsSet("file") {
		log.Fatal("Missing file parameter!")
	}
	backendInstance.Distribute(viper.GetString("file"))
}
