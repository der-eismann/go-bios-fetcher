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

package asus

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
	"github.com/sirupsen/logrus"
)

func getHashID(website string) (string, error) {
	// Downloading the given website
	response, err := http.Get(website)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Saving the website as a string
	bodyRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyRead)

	// Getting the model info JSON within the website
	regex := regexp.MustCompile(`(?sU)pdhashid: "(?P<content>.*)"`)
	result := regex.FindStringSubmatch(bodyString)
	var content string
	if len(result) > 1 {
		content = result[1]
	} else {
		return "", errors.New("Regex couldn't find hash ID")
	}

	return content, nil
}

func getDownloadJSON(hashID string) ([]byte, error) {
	// Downloading the given website
	response, err := http.Get("https://www.asus.com/support/api/product.asmx/GetPDDrivers?osid=45&website=en&pdhashedid=" + hashID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Saving the website as a string
	bodyRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyRead)

	return []byte(bodyString), nil
}

// func GetMockJSON() []byte {
// 	file, _ := os.Open("downloads.json")
// 	content, _ := ioutil.ReadAll(file)
// 	return content
// }

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {
	hashID, err := getHashID(device.URL)
	if err != nil {
		logrus.Error(err)
	}
	logrus.Debugf("Hash ID: %s", hashID)

	downloadJSON, err := getDownloadJSON(hashID)
	parsedDownloads := downloadsASUS{}
	err = json.Unmarshal(downloadJSON, &parsedDownloads)
	if err != nil {
		logrus.Error(err)
	}

Loop:
	for filterPos, filter := range device.Downloads {
		for _, category := range parsedDownloads.Result.Objects {
			for _, download := range category.Files {
				if strings.Contains(download.Title, filter.Filter) {
					device.Downloads[filterPos].Version = download.Version
					device.Downloads[filterPos].Link = download.DownloadURL.Global
					releaseDate, _ := time.Parse("2006/01/02", download.ReleaseDate)
					device.Downloads[filterPos].Date = releaseDate.Format("02.01.2006")
					device.Downloads[filterPos].Readme = download.Description
					continue Loop
				}
			}
		}
	}
	logrus.Debugf("%#v", device)

	return device
}
