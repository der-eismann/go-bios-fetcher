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
	"strings"
	"time"

	"github.com/der-eismann/go-bios-fetcher/vendors/lenovo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all entries defined in config.toml",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func update() {
	currentData := viper.Get("devices")

	for i, devicesInterface := range currentData.([]interface{}) {
		device := devicesInterface.(map[string]interface{})
		newVersion := ""
		newFile := ""

		// Only update devices that are enabled
		enabled := !(device["enabled"].(int64) == 0)
		if enabled {
			switch device["vendor"] {
			case "Lenovo":
				newVersion, newFile = lenovo.GetLatestBios(device["url"].(string))
			default:
				fmt.Println("No compatible vendor defined!")
				continue
			}

			if _, ok := device["version"]; ok {
				if strings.Compare(device["version"].(string), newVersion) != 0 {
					fmt.Printf("Neue Version für %s: Von %s auf %s aktualisiert!\n", device["model"], device["version"], newVersion)
				} else {
					fmt.Printf("%s: BIOS %s ist auf dem aktuellen Stand!\n", device["model"], device["version"])
				}
			} else {
				fmt.Printf("Aktuelles BIOS für %s: Version %s zum Download verfügbar!\n", device["model"], newVersion)
			}

			currentData.([]interface{})[i].(map[string]interface{})["version"] = newVersion
			currentData.([]interface{})[i].(map[string]interface{})["download"] = newFile
		}
	}

	viper.Set("date", time.Now())
	err := viper.WriteConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error writing config file: %s \nConfig file: %s\n", err, viper.ConfigFileUsed()))
	}
}
