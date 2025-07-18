package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"monitoring/checks"
	"monitoring/types"
)

func main() {
	cfgFile := "config/config.yaml"
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		os.Exit(1)
	}

	result := checks.RunAllChecks(cfg)

	outputPath := cfg.General.OutputPath
	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		fmt.Printf("Failed to create output dir: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputPath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Wrote output to %s\n", outputPath)
}
