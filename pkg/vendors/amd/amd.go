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

package amd

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
	"github.com/sirupsen/logrus"
)

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {
	client := &http.Client{}

	req, err := http.NewRequest("GET", device.URL, nil)
	if err != nil {
		logrus.Error(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error(err)
	}

	regex := regexp.MustCompile(`(?sU)AMD Chipset Drivers.*<div class="field__item">(.*)<\/div>.*00Z">(\d\d\.\d\d\.\d\d\d\d).*<a href="(.*)"`)
	result := regex.FindStringSubmatch(string(body))

	if len(result) > 3 {
		device.Downloads[0].Version = result[1]
		device.Downloads[0].Date = result[2]
		device.Downloads[0].Link = result[3]
	} else {
		logrus.Error(errors.New("Regex was not successful"))
	}
	logrus.Printf("%#v", device)

	return device
}
