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
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
)

func GetWebsiteContent(website string) ([]byte, error) {
	logrus.Debugf("Downloading %s", website)
	// Downloading the given website
	cookieJar, _ := cookiejar.New(nil)
	httpClient := http.Client{
		Jar: cookieJar,
	}
	response, err := httpClient.Get(website)
	cmdutil.Must(err)
	defer response.Body.Close()

	// Saving the website as a byte array
	bodyRead, err := ioutil.ReadAll(response.Body)
	cmdutil.Must(err)

	return bodyRead, nil
}

// GetDownloadJSON is downloading all the information about the available
// downloads from the given Lenovo support website
func GetDownloadJSON(website string) ([]byte, error) {
	downloadPage, err := GetWebsiteContent(website)
	cmdutil.Must(err)

	// Parsing the productID within the website
	regex := regexp.MustCompile(`(?m)var config = window\.config \|\| (.*);`)
	result := regex.FindStringSubmatch(string(downloadPage))
	windowConfigJSON := windowConfig{}
	json.Unmarshal([]byte(result[1]), &windowConfigJSON)

	actualDownloads, err := GetWebsiteContent("https://pcsupport.lenovo.com/de/de/api/v4/downloads/drivers?productId=" + windowConfigJSON.DynamicItems.ProductID)
	cmdutil.Must(err)
	logrus.Debugf("%s", actualDownloads)

	return actualDownloads, nil
}

func GetMockJSON() []byte {
	file, _ := os.Open("drivers.json")
	content, _ := ioutil.ReadAll(file)
	return content
}

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {

	decodedData, err := GetDownloadJSON(device.URL)
	if err != nil {
		logrus.Println(err)
	}
	//decodedData := GetMockJSON()

	parsedJSON := downloads{}
	json.Unmarshal(decodedData, &parsedJSON)

	for filterPos, filter := range device.Downloads {
		for _, download := range parsedJSON.Body.DownloadItems {
			if strings.Contains(download.Title, filter.Filter) {
				for _, file := range download.Files {
					logrus.Debugf("%#v", file)
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
