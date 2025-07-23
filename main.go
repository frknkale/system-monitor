package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"monitoring/checks"
	"monitoring/logger"
	"monitoring/types"
)

func main() {
	cfgFile := "config/config.yaml"

	// Read config file
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Printf("Failed to parse config file: %v\n", err)
		logger.Log.Fatalf("Failed to parse config file: %v", err)
		os.Exit(1)
	}
	
	logPath := cfg.General.LogPath
	outputPath := cfg.General.OutputPath

	// Initialize logger
	err = logger.Init(logPath)
	if err != nil {
		fmt.Printf("Logger initialization failed: %v\n", err)
		os.Exit(1)
	}
	logger.Log.Println("Logger initialized.")


	// Determine interval
	intervalStr := cfg.General.Interval
	interval, err := time.ParseDuration(intervalStr)
	if err != nil || interval <= 0 || interval > time.Hour{
		fmt.Printf("Invalid interval format (%s), defaulting to 30s\n", intervalStr)
		interval = 30 * time.Second
	}

	for {
		logger.Log.Println("Running system checks...")
		result := checks.RunAllChecks(cfg)

		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			logger.Log.Printf("Failed to create output dir: %v", err)
			fmt.Printf("Failed to create output dir: %v", err)
			continue
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			logger.Log.Printf("Failed to marshal JSON: %v", err)
			continue
		}

		err = os.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			logger.Log.Printf("Failed to write output file: %v", err)
			fmt.Printf("Failed to write output file: %v", err)
			continue
		}

		logger.Log.Println("System check completed.")
		fmt.Println("System check completed.")
		fmt.Println("Wrote output to %s\n", outputPath)

		logger.Log.Println(string(jsonData))
		
		fmt.Println(">>> Sleeping for:", interval)

		time.Sleep(interval)

	}
}
