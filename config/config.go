package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	OxygenURL string     `yaml:"oxygen_url"`
	LogFile   string     `yaml:"log_file"`
	Schedules []Schedule `yaml:"schedules"`
}

type Schedule struct {
	Hour   int `yaml:"hour"`
	Minute int `yaml:"minute"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		exampleData, exampleErr := ioutil.ReadFile(filename + ".example")
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
