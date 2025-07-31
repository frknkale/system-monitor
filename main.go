package main

import (
	"monitoring/cache"
	"monitoring/config"
	"monitoring/monitoring"
	"monitoring/webserver"
	"monitoring/alerter"
)

func main() {
	cache.SetCache(map[string]interface{}{"status": "loading"})

	config.ReadConfig("config/config.yaml")
	cfg:= config.GetConfig()
	
	alerter.Init(cfg)
	
	go monitoring.Monitoring(cfg)
	go webserver.WebServer(cfg)
	select {}
}