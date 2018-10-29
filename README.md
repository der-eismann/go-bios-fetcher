# go-bios-fetcher
Simple tool to fetch the latest BIOS version from multiple hardware vendors

[![Go Report Card](https://goreportcard.com/badge/github.com/der-eismann/go-bios-fetcher)](https://goreportcard.com/report/github.com/der-eismann/go-bios-fetcher) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Motivation

As a former sysadmin with an affection for BIOS updates I wanted to keep track of them without RSS feeds, newsletters or whatsoever. At my last workplace I started designing a PHP script for this following a 'works for me' concept with a PostgreSQL database behind. With the aim of learning Go I wanted to port the project and at the same time make it more portable, extensible and open source. 


## Goals

 * Support all major hardware vendors
 * Make it as easy to use as possible
 * Extensibility
 * A+ Rating on [Report Card](https://goreportcard.com/badge/github.com/der-eismann/go-bios-fetcher)

## Progress

### Supported Vendors
 * Lenovo

### Planned Vendors
 * ASUS
 * Gigabyte

### "If the need arises" Vendors
 * Acer
 * ASRock
 * MSI
 * Supermicro

## Installation

1. Clone repo and enter folder
2. Fetch dependencies
3. Run `go build`

## Configuration

The configuration called `config.toml` must be located in the same folder as the `go-bios-fetcher` binary.

Example:
``` toml
date = 2018-10-30T15:00:00+01:00

[[devices]]
  enabled = 1
  vendor = "Lenovo"
  model = "ThinkPad T480s"
  url = "https://pcsupport.lenovo.com/de/de/products/laptops-and-netbooks/thinkpad-t-series-laptops/thinkpad-t480s-type-20l7-20l8/downloads"
  version = "1.26"
  download = "https://download.lenovo.com/pccbbs/mobiles/n22uj10w.exe"

[[devices]]
  enabled = 0 // Exclude entry from being updated
  vendor = "Lenovo"
  model = "ThinkPad T470"
  url = "https://pcsupport.lenovo.com/de/de/products/laptops-and-netbooks/thinkpad-t-series-laptops/thinkpad-t470/downloads"
```

## Usage

After defining all your devices in the config file you can use the following commands:

| CMD      | Description                                                    |
| -------- | -------------------------------------------------------------- |
| help     | Get help for the other commands.                               |
| list     | List the locally stored information about the latest versions. |
| update   | Fetch new BIOS versions for all defined & enabled devices.     |

## License
 
Copyright 2018 Philipp Trulson <philipp@trulson.de>
 
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
 
http://www.apache.org/licenses/LICENSE-2.0
 
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
