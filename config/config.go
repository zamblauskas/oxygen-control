package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OxygenURL string      `yaml:"oxygen_url"`
	LogFile   string      `yaml:"log_file"`
	DBPath    string      `yaml:"db_path"`
	Schedule  *Schedule   `yaml:"schedule,omitempty"`
	Flic      *FlicConfig `yaml:"flic,omitempty"`
}

type Schedule struct {
	Enabled bool `yaml:"enabled"`
}

type FlicConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ServerURL string `yaml:"server_url"`
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
