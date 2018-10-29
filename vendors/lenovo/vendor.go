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

package lenovo

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/pierrec/lz4"
)

// GetLatestBios polls the given website and returns version & link to file
func GetLatestBios(website string) (string, string) {
	// Variable for saving the array index of BIOS download
	var arrayKey string
	// Variables for the return values
	var version, url string

	// Downloading the given website
	response, err := http.Get(website)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer response.Body.Close()

	// Saving the website as a string
	bodyRead, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	bodyString := string(bodyRead)

	// Parsing the encoded JSON within the website
	regex := regexp.MustCompile(`ds_downloads.*content":"(?P<content>.*)","originLength":(?P<length>.*)}`)
	result := regex.FindStringSubmatch(bodyString)
	content := result[1]
	length, err := strconv.Atoi(result[2])
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Decoding the JSON from Base64
	lz4Data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// Decoding the JSON from LZ4
	decodedData := make([]byte, length)
	_, err = lz4.UncompressBlock(lz4Data, decodedData)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	arrayIndex := 0
	// Looping through every download to find the one containing the BIOS update
	jsonparser.ArrayEach(decodedData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		downloadTitle, err := jsonparser.GetString(value, "Title")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if strings.Contains(downloadTitle, "BIOS Update Utility") {
			// Saving the array index in the form of [21]
			arrayKey = fmt.Sprintf("[%d]", arrayIndex)
		}
		arrayIndex++
	}, "DownloadItems")

	// Looping through every file for the BIOS update, skipping Readmes and ISOs
	jsonparser.ArrayEach(decodedData, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		versionTemp, err := jsonparser.GetString(value, "Version")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		typeString, err := jsonparser.GetString(value, "TypeString")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		urlTemp, err := jsonparser.GetString(value, "URL")
		if err != nil {
			fmt.Println("Error: ", err)
		}

		if strings.Contains(typeString, "EXE") {
			version = versionTemp
			url = urlTemp
		}
	}, "DownloadItems", arrayKey, "Files")

	return version, url
}
