package main

import (
	"monitoring/cache"
	"monitoring/config"
	"monitoring/monitoring"
	"monitoring/webserver"
	"monitoring/alerter"
	"os"
	"fmt"
)

func main() {
	configPath := "/opt/monitoring/config/config.yaml"

	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}
	fmt.Printf("Using config file: %s\n", configPath)
	cache.SetCache(map[string]interface{}{"status": "loading"})

	config.ReadConfig(configPath)
	cfg:= config.GetConfig()

	alerter.Init(cfg)
	
	go monitoring.Monitoring(cfg)
	go webserver.WebServer(cfg)
	select {}
}