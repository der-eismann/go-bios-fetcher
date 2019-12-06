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
	"io/ioutil"
	"time"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
	"github.com/der-eismann/go-bios-fetcher/pkg/vendors/asus"
	"github.com/der-eismann/go-bios-fetcher/pkg/vendors/lenovo"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func update() {
	config, err := lib.ReadConfig()
	cmdutil.Must(err)

	for device := range config.Devices {
		logrus.Printf("Downloading files for %s...", config.Devices[device].Name)
		switch config.Devices[device].Vendor {
		case "Lenovo":
			config.Devices[device] = lenovo.GetLatestFiles(config.Devices[device])
		case "Asus":
			config.Devices[device] = asus.GetLatestFiles(config.Devices[device])
		}

	}
	config.LastUpdated = time.Now().Format("02.01.2006 15:04")

	marshalled, err := yaml.Marshal(&config)
	cmdutil.Must(err)
	err = ioutil.WriteFile("config.yaml", marshalled, 0644)
	cmdutil.Must(err)
}

type UpdateApp struct {
}

func (app *UpdateApp) Run(ctx context.Context, cmd *cobra.Command, args []string) {
	update()
}

func NewUpdateCommand() *cobra.Command {
	app := new(UpdateApp)
	cmd := cmdutil.NewCommand(app)
	cmd.Use = "update"
	cmd.Short = "Updates all entries in the config.yaml"
	return cmd
}
