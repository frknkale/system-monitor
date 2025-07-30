package webserver

import (
	"encoding/json"
	"net/http"
	"time"
	"os"
	"fmt"
	"bufio"
	// "bytes"
	// "io/ioutil"
)

type Status struct {
	Timestamp string `json:"time"`
	Host      string `json:"host"`
	Disk      map[string]interface{} `json:"disks"`
	Memory    map[string]interface{} `json:"memory"`
	CPU       map[string]interface{} `json:"cpu"`
	Services  map[string]interface{} `json:"services"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	path := "/var/log/web/test-server-output.json"

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
		data, err := readLatestStatus("/var/log/web/test-server-output.json")
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


func WebServer() {
	http.HandleFunc("/healthz", healthHandler)
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


