// Copyright © 2019 Philipp Trulson <philipp@trulson.de>
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
	"context"
	"fmt"
	"strings"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/spf13/cobra"
)

func list(showUrl bool) {
	config, err := lib.ReadConfig()
	cmdutil.Must(err)

	fmt.Printf("Last updated at %s\n\n", config.LastUpdated)
	for _, device := range config.Devices {
		fmt.Println(device.Name)
		var lengthName, lengthVersion int
		for _, download := range device.Downloads {
			if len(download.Filter) > lengthName {
				lengthName = len(download.Filter)
			}
			if len(download.Version) > lengthVersion {
				lengthVersion = len(download.Version)
			}
		}
		for _, download := range device.Downloads {
			fmt.Printf("%-[2]*[1]s | %-[4]*[3]s | %[5]s\n", download.Filter, lengthName, download.Version, lengthVersion, download.Date)
			if showUrl {
				fmt.Printf("  - Download: %s\n", download.Link)
				fmt.Printf("  - Readme:   %s\n\n", strings.Replace(download.Link, "exe", "txt", 1))
			}
		}
		fmt.Println("")
	}
}

type ListApp struct {
	ShowUrl bool
}

func (app *ListApp) Run(ctx context.Context, cmd *cobra.Command, args []string) {
	list(app.ShowUrl)
}

func (app *ListApp) Bind(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(
		&app.ShowUrl, "show-url", false,
		`Set to true to show the full download URLs.`)
}

func NewListCommand() *cobra.Command {
	app := new(ListApp)
	cmd := cmdutil.NewCommand(app)
	cmd.Use = "list"
	cmd.Short = "Lists all entries in the config.yaml"
	return cmd
}
