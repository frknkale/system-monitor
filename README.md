# Go Based System Monitor

A Go-based system monitoring application that collects metrics like CPU, memory, disk usage, processes, permissions, services and network stats. It can send metrics to an ELK stack, log to local files, alert if thresholds exceed, and serve a web interface with real-time data.

---

## Features

- Monitor CPU, memory, disk, network, permission, processes, and services
- Configurable via YAML config file
- Send metrics remotely (e.g., to ELK stack)
- Web server to view metrics dashboard
- Alerting based on customizable thresholds
- Logs metrics and alerts locally

---

## Configuration

Everything is configured via the `config.yaml` file. Adjust the settings there to customize metrics collection, output paths, remote server settings, etc.

- If you want to send the collected metrics output to a remote server for tools like ELK stack, enable the remote option (`general.remote.enabled`) in the configuration file. Make sure to configure the remote host, user, and path accordingly.

---

## Running the Application

Make sure you have root privileges to allow the app to access system metrics.

Run directly with:
```bash
go run main.go /path/to/config.yaml
```

or build and run the binary:

```bash
go build -o monitoring-app
./monitoring-app /path/to/config.yaml
```
The monitoring web server will start on port `:8080` by default and serve the metrics dashboard.

---

## Note

The metrics output is cached and updated at configured intervals (e.g., every 30 seconds). If you want to get fresh, immediately updated data from the web server, you can append ```?fresh=true``` to the URL. This forces the app to refresh the data.

Example:

```
http://localhost:8080/?fresh=true
```
