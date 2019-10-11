package lib

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	LastUpdated string   `yaml:"lastUpdated`
	Devices     []Device `yaml:"devices"`
}

type Device struct {
	Name      string     `yaml:"name"`
	Vendor    string     `yaml:"vendor"`
	URL       string     `yaml:"url"`
	Downloads []Download `yaml:"downloads"`
}

type Download struct {
	Filter  string `yaml:"filter"`
	Version string `yaml:"version"`
	Date    string `yaml:"date"`
	Link    string `yaml:"link"`
	Readme  string `yaml:"readme"`
}

func ReadConfig() (Config, error) {
	config := Config{}
	file, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, err
	}
	configContent, err := ioutil.ReadAll(file)
	if err != nil {
		return Config{}, err
	}
	err = yaml.Unmarshal(configContent, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
