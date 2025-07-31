package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"monitoring/checks"
	"monitoring/logger"
	"monitoring/types"
	"monitoring/cache"
)

func Monitoring(cfg types.Config) {
	// Read config file
	
	logPath := cfg.General.LogPath
	outputPath := cfg.General.OutputPath
	alertPath := cfg.Alerter.LogPath
	remote := cfg.General.Remote

	// Initialize logger
	err := logger.Init(logPath)
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

		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		if err != nil {
			logger.Log.Printf("Failed to create output dir: %v", err)
			fmt.Printf("Failed to create output dir: %v", err)
			continue
		}

		result := checks.RunAllChecks(cfg)

		cache.SetCache(result)		// Store the result in the shared cache for Web Server

		jsonData, err := json.Marshal(result)
		if err != nil {
		    logger.Log.Printf("Failed to marshal JSON: %v", err)
		    continue
		}

		f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
		    logger.Log.Printf("Failed to open output file: %v", err)
		    continue
		}
		defer f.Close()

		if _, err := f.Write(jsonData); err != nil {
		    logger.Log.Printf("Failed to write JSON to file: %v", err)
		    continue
		}
		if _, err := f.Write([]byte("\n")); err != nil {
		    logger.Log.Printf("Failed to write newline to file: %v", err)
		    continue
		}

		if remote.Enabled{
			user, host, path:= remote.User, remote.Host, remote.RemotePath
			rsyncCmd := exec.Command(
				"sudo", "-u", user, "rsync", "--inplace", "-az", outputPath,
				fmt.Sprintf("%s@%s:%s", user, host, path))
				
			if err := rsyncCmd.Run(); err != nil {
				logger.Log.Printf("Failed to rsync output.json to %s@%s:%s: %v",user, host, path, err)
				fmt.Printf("Failed to rsync output.json to %s@%s:%s: %v\n",user, host, path, err)
			} else {
				logger.Log.Printf("Successfully rsynced output.json to %s@%s:%s",user, host, path)
				fmt.Printf("Successfully rsynced output.json to %s@%s:%s\n",user, host, path)
			}
		}

		if cfg.Alerter.Enabled {
			user, host, path := cfg.Alerter.Remote.User, cfg.Alerter.Remote.Host, cfg.Alerter.Remote.RemotePath
			rsyncCmd := exec.Command(
				"sudo", "-u", user, "rsync", "--inplace", "-az", alertPath,
				fmt.Sprintf("%s@%s:%s", user, host, path))
			
			if err := rsyncCmd.Run(); err != nil {
				logger.Log.Printf("Failed to rsync alerts.json to %s@%s:%s: %v", user, host, path, err)
				fmt.Printf("Failed to rsync alerts.json to %s@%s:%s: %v\n", user, host, path, err)
			} else {
				logger.Log.Printf("Successfully rsynced alerts.json to %s@%s:%s", user, host, path)
				fmt.Printf("Successfully rsynced alerts.json to %s@%s:%s\n", user, host, path)
			}
		}
		

		logger.Log.Println("Checks completed.")
		fmt.Println("Checks completed.")
		fmt.Printf("Wrote output to %s\n", outputPath)
		
		// fmt.Println(string(jsonData))

		logger.Log.Println(string(jsonData))
		
		fmt.Println(">>> Sleeping for:", interval)

		time.Sleep(interval)
	}
}


