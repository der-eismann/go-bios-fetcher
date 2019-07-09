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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pierrec/lz4"
	"github.com/sirupsen/logrus"

	"github.com/der-eismann/go-bios-fetcher/pkg/lib"
)

// GetDownloadJSON is downloading all the information about the available
// downloads from the given Lenovo support website
func GetDownloadJSON(website string) ([]byte, error) {
	// Downloading the given website
	response, err := http.Get(website)
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

	// Parsing the encoded JSON within the website
	regex := regexp.MustCompile(`ds_downloads.*content":"(?P<content>.*)","originLength":(?P<length>.*)}`)
	result := regex.FindStringSubmatch(bodyString)
	content := result[1]
	length, err := strconv.Atoi(result[2])
	if err != nil {
		return nil, err
	}

	// Decoding from Base64 to LZ4 encoded data
	lz4Data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	// Decoding the JSON from LZ4
	decodedData := make([]byte, length)
	_, err = lz4.UncompressBlock(lz4Data, decodedData)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

func GetMockJSON() []byte {
	file, _ := os.Open("downloads.json")
	content, _ := ioutil.ReadAll(file)
	return content
}

// GetLatestFiles polls the given website and returns version & link to file
func GetLatestFiles(device lib.Device) lib.Device {

	decodedData, _ := GetDownloadJSON(device.URL)
	//decodedData := GetMockJSON()

	parsedJSON := downloads{}
	json.Unmarshal(decodedData, &parsedJSON)

	for filterPos, filter := range device.Downloads {
		for _, download := range parsedJSON.DownloadItems {
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
						break
					}
				}
			}
		}
	}

	return device
}
