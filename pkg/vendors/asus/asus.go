//https://www.asus.com/support/api/product.asmx/GetPDDrivers?cpu=&osid=45&website=de&pdhashedid=xlQ4iKfiVURLAgq3&model=ROG%20STRIX%20X570-E%20GAMING&callback=supportpdpage
//https://www.asus.com/support/api/product.asmx/GetPDBIOS?website=de&pdhashedid=xlQ4iKfiVURLAgq3&model=ROG%20STRIX%20X570-E%20GAMING&cpu=&callback=supportpdpage
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
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
	"github.com/sirupsen/logrus"
)

// GetHashID gets the HashID of the given product
func GetHashID(website string) (string, error) {
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
	content := result[1]

	return content, nil
}

// GetDownloadJSON is downloading all the information about the available
// downloads from the given Lenovo support website
func GetDownloadJSON(hashID string) ([]byte, error) {
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

func GetMockJSON() []byte {
	file, _ := os.Open("downloads.json")
	content, _ := ioutil.ReadAll(file)
	return content
}

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {

	hashID, err := GetHashID(device.URL)
	if err != nil {
		logrus.Error(err)
	}
	downloadJSON, err := GetDownloadJSON(hashID)
	parsedDownloads := downloadsASUS{}
	err = json.Unmarshal(downloadJSON, &parsedDownloads)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Debugf("%#v", parsedDownloads)

Loop:
	for filterPos, filter := range device.Downloads {
		for _, category := range parsedDownloads.Result.Objects {
			for _, download := range category.Files {
				if strings.Contains(download.Title, filter.Filter) {
					device.Downloads[filterPos].Version = download.Version
					device.Downloads[filterPos].Link = download.DownloadURL.Global
					device.Downloads[filterPos].Date = download.ReleaseDate
					continue Loop
				}
			}
		}
	}

	return device
}
