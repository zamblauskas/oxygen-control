package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OxygenURL string      `yaml:"oxygen_url"`
	LogFile   string      `yaml:"log_file"`
	Schedules []Schedule  `yaml:"schedules"`
	Flic      *FlicConfig `yaml:"flic,omitempty"`
}

type Schedule struct {
	Hour   int `yaml:"hour"`
	Minute int `yaml:"minute"`
}

type FlicConfig struct {
	Enabled             bool   `yaml:"enabled"`
	ServerURL           string `yaml:"server_url"`
	ButtonBluetoothAddr string `yaml:"button_bluetooth_address"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		exampleData, exampleErr := os.ReadFile(filename + ".example")
		if exampleErr != nil {
			return nil, err
		}
		data = exampleData
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
