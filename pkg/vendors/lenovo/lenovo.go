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

package lenovo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
)

func getWebsiteContent(website string) ([]byte, error) {
	logrus.Debugf("Downloading %s", website)
	// Creating a cookie jar - the Lenovo website needs cookies
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	httpClient := http.Client{
		Jar: cookieJar,
	}

	// Downloading the given website
	response, err := httpClient.Get(website)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}
	defer response.Body.Close()

	// Saving the website as a byte array
	bodyRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	return bodyRead, nil
}

func getDownloadJSON(website string) ([]byte, error) {
	// Load the support site of the device
	downloadPage, err := getWebsiteContent(website)
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	// Parsing the productID within the website
	regex := regexp.MustCompile(`(?m)var config = window\.config \|\| (.*);`)
	result := regex.FindStringSubmatch(string(downloadPage))
	windowConfigJSON := windowConfig{}
	if len(result) > 1 {
		err = json.Unmarshal([]byte(result[1]), &windowConfigJSON)
		if err != nil {
			return []byte{}, errors.WithStack(err)
		}
	} else {
		return []byte{}, errors.New("Regex couldn't find product ID")
	}

	// Dowload the JSON file with the actual downloads
	logrus.Debugf("Using product ID %s for downloads...", windowConfigJSON.DynamicItems.ProductID)
	actualDownloads, err := getWebsiteContent(fmt.Sprintf("https://pcsupport.lenovo.com/de/de/api/v4/downloads/drivers?productId=%s", windowConfigJSON.DynamicItems.ProductID))
	if err != nil {
		return []byte{}, errors.WithStack(err)
	}

	return actualDownloads, nil
}

func getMockJSON() []byte {
	file, _ := os.Open("drivers.json")
	content, _ := ioutil.ReadAll(file)
	return content
}

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {

	decodedData, err := getDownloadJSON(device.URL)
	if err != nil {
		logrus.Error(err)
	}
	//decodedData := getMockJSON()

	parsedJSON := downloads{}
	err = json.Unmarshal(decodedData, &parsedJSON)
	if err != nil {
		logrus.Error(err)
	}

	for filterPos, filter := range device.Downloads {
		for _, download := range parsedJSON.Body.DownloadItems {
			if strings.Contains(download.Title, filter.Filter) {
				for _, file := range download.Files {
					if file.TypeString == "EXE" || file.TypeString == "zip" {
						if strings.Contains(filter.Filter, "BIOS") && (strings.Contains(file.Name, "setting") || file.TypeString == "zip") {
							continue
						}
						device.Downloads[filterPos].Version = file.Version
						device.Downloads[filterPos].Link = file.URL
						device.Downloads[filterPos].Date = time.Unix(file.Date.Unix/1000, 0).Format("02.01.2006")
					}
					if file.TypeString == "TXT README" {
						device.Downloads[filterPos].Readme = file.URL
					}
				}
			}
		}
	}
	logrus.Debugf("%#v", device)

	return device
}
