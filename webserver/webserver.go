package webserver

import (
	"encoding/json"
	"net/http"
	"html/template"
	"time"
	"os"
	"fmt"
	"bufio"
    "strconv"
    "strings"
	"monitoring/types"
	// "bytes"
	// "io/ioutil"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/log/monitoring/metrics/output.json"

	file, err := os.Open(path)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to open file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lastLine = line
		}
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, fmt.Sprintf("Failed to scan file: %v", err), http.StatusInternalServerError)
		return
	}

	if lastLine == "" {
		http.Error(w, "No valid lines in status file", http.StatusInternalServerError)
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(lastLine), &result); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func readLatestStatus(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			lastLine = line
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal([]byte(lastLine), &result); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return result, nil
}

func sectionHandler(section string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := readLatestStatus("/var/log/monitoring/metrics/output.json")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sectionData, ok := data[section]
		if !ok {
			http.Error(w, fmt.Sprintf("Section '%s' not found", section), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sectionData)
	}
}


func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	raw, err := readLatestStatus("/var/log/monitoring/metrics/output.json")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var data types.DashboardData

	// Timestamp
	if ts, ok := raw["timestamp"].(string); ok {
		data.Timestamp = ts
	}

	// CPU Usage
	if cpuMap, ok := raw["cpu"].(map[string]interface{}); ok {
		if val, ok := cpuMap["usage_percent"].(string); ok {
		    val = strings.TrimSuffix(val, "%")
		    if f, err := strconv.ParseFloat(val, 64); err == nil {
		        data.CPUPercent = f
		    }
		}
	}

	// Memory Usage (string percent -> float)
	if memMap, ok := raw["memory"].(map[string]interface{}); ok {
		if strVal, ok := memMap["used_percent"].(string); ok {
			strVal = strings.TrimSuffix(strVal, "%")
			if fval, err := strconv.ParseFloat(strVal, 64); err == nil {
				data.MemoryPercent = fval
			}
		}
	}

	// Disk Usage (first disk, first partition)
	if diskArray, ok := raw["disk"].([]interface{}); ok && len(diskArray) > 0 {
		if diskEntry, ok := diskArray[0].(map[string]interface{}); ok {
			if parts, ok := diskEntry["partitions"].([]interface{}); ok {
				for _, part := range parts {
					if partMap, ok := part.(map[string]interface{}); ok {
						if mount, ok := partMap["mountpoint"].(string); ok && mount == "/" {
							if val, ok := partMap["used_percent"].(float64); ok {
								data.DiskRootPercent = val
							}
						}
					}
				}
			}
		}
	}

	tmpl, err := template.ParseFiles("webserver/templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl.Execute(w, data)
}


func WebServer() {
	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/status", statusHandler)

	http.HandleFunc("/cpu", sectionHandler("cpu"))
	http.HandleFunc("/memory", sectionHandler("memory"))
	http.HandleFunc("/disk", sectionHandler("disk"))
	http.HandleFunc("/network", sectionHandler("network"))
	http.HandleFunc("/services", sectionHandler("services"))
	http.HandleFunc("/permissions", sectionHandler("permissions"))
	http.HandleFunc("/processes", sectionHandler("processes"))

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}


