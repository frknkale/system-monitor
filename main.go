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
		os.Exit(1)
	}

	// Initialize logger
	logPath := "logs/monitor.log"
	if logPath == "" {
		logPath = "logs/monitor.log"
	}
	err = logger.Init(logPath)
	if err != nil {
		fmt.Printf("Logger initialization failed: %v\n", err)
		os.Exit(1)
	}
	logger.Log.Println("Logger initialized.")

	// Determine interval
	interval := time.Duration(cfg.General.IntervalSeconds) * time.Second
	if interval <= 0 || interval > time.Hour {
		fmt.Printf("Invalid interval (%d), defaulting to 30s\n", cfg.General.IntervalSeconds)
		interval = 30 * time.Second
	}

	for {
		logger.Log.Println("Running system checks...")
		result := checks.RunAllChecks(cfg)

		outputPath := cfg.General.OutputPath
		if outputPath == "" {
			outputPath = "output/output.json"
		}

		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			logger.Log.Printf("Failed to create output dir: %v", err)
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
			continue
		}

		logger.Log.Println("System check complete.")
		fmt.Println(string(jsonData))
		fmt.Println(">>> Sleeping for:", interval)
		time.Sleep(interval)
		// fmt.Println(">>> Awake from sleep")
	}
}
