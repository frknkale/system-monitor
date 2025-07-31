package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"

	"monitoring/types"
)

var config types.Config

func ReadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}
}

func GetConfig() types.Config {
	return config
}