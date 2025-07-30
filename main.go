package main

import (
	"monitoring/monitoring"
	"monitoring/webserver"
)

func main() {
	go monitoring.Monitoring("config/config.yaml")
	go webserver.WebServer()
	select {}
}