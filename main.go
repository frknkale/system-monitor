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
	// Initialize logger
	logPath := "logs/monitor.log"
	err := logger.Init(logPath)
	if err != nil {
		fmt.Printf("Logger initialization failed: %v\n", err)
		os.Exit(1)
	}
	logger.Log.Println("Logger initialized.")

	cfgFile := "config/config.yaml"
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		logger.Log.Fatalf("Failed to read config file: %v", err)
	}

	var cfg types.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		logger.Log.Fatalf("Failed to parse config file: %v", err)
	}

	interval := time.Duration(cfg.General.IntervalSeconds) * time.Second
	if interval == 0 {
		interval = 30 * time.Second // default fallback
	}

	// Continuous monitoring loop
	for {
		logger.Log.Println("Running system checks...")
		result := checks.RunAllChecks(cfg)

		outputPath := cfg.General.OutputPath
		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			logger.Log.Fatalf("Failed to create output dir: %v", err)
		}
		fmt.Printf("Wrote output to %s\n", outputPath)

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			logger.Log.Fatalf("Failed to marshal JSON: %v", err)
		}

		err = os.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			logger.Log.Fatalf("Failed to write output file: %v", err)
		}

		logger.Log.Println(string(jsonData)) // log full result
		// fmt.Printf(string(jsonData))
		time.Sleep(interval)
	}
}
