package main

import (
	//	"fmt"
	//	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port     string   `yaml:"port"`
	Backends []string `yaml:"backends"`
	Timeouts Timeouts `yaml:"timeouts"`
}

type Timeouts struct {
	Read string `yaml:"read"`
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
