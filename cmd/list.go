// Copyright © 2018 Philipp Trulson <philipp@trulson.de>
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
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show all entries with their latest version",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		list()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func list() {
	currentData := viper.Get("devices")
	lastUpdated := viper.Get("date").(time.Time)

	for _, devicesInterface := range currentData.([]interface{}) {
		device := devicesInterface.(map[string]interface{})
		if device["enabled"] != int64(1) {
			continue
		}
		fmt.Printf("%s %s:\n", device["vendor"], device["model"])
		fmt.Printf("  - BIOS: %s\n", device["version"])
		fmt.Printf("  - Download: %s\n\n", device["download"])
	}
	fmt.Printf("Last updated: %s\n", lastUpdated.Format("02.01.2006 15:04"))
}
