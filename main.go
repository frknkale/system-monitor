package main

import (
	"monitoring/monitoring"
	"monitoring/webserver"
	"monitoring/cache"
)

func main() {
	cache.SetCache(map[string]interface{}{"status": "loading"})

	err := webserver.LoadConfig("config/config.yaml")
	if err != nil {
		panic(err)
	}
	
	go monitoring.Monitoring("config/config.yaml")
	go webserver.WebServer()
	select {}
}