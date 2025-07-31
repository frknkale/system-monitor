package webserver

import (
	"encoding/json"
	"net/http"
	"html/template"
	"time"
	"os"
	"fmt"
    "strconv"
    "strings"
	"gopkg.in/yaml.v3"

	"monitoring/types"
	"monitoring/checks"
	"monitoring/cache"
	"monitoring/checks/systemchecks"
	
	// "bytes"
	// "io/ioutil"
)

var (
	config    types.Config
	sem = make(chan struct{}, 1)
	i = 0
)

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config)
}

func getData(fresh bool) map[string]interface{} { // If fresh is true, or the cache is expired, fetch fresh data
	data := cache.GetCache()
	if len(sem) == 0 && (fresh || cache.IsExpired()) {
		sem <- struct{}{}
		go func() {
			defer func() {
				<-sem
			}()
			i++
			fmt.Println("Refreshing Data...")
			fmt.Println("Number of goroutines running concurrently:",i)
			data := checks.RunAllChecks(config)
			i--
			cache.SetCache(data)
			fmt.Println("Refresh ended.")
		}()
	}

	return data  // Serve the old cache immediately while the fresh data is being fetched
}


func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	fresh := r.URL.Query().Get("fresh") == "true"   // Check for ?fresh=true

	data:= getData(fresh)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func sectionHandler(section string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fresh := r.URL.Query().Get("fresh") == "true"   // Check for ?fresh=true

		data:= getData(fresh)

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
	fresh := r.URL.Query().Get("fresh") == "true"   // Check for ?fresh=true

	data:= getData(fresh)

	var dashboardData types.DashboardData

	dashboardData.Timestamp = cache.GetLastUpdated().Format(time.RFC3339)

	// CPU Usage
	if cpuMap, ok := data["cpu"].(map[string]interface{}); ok {
		if val, ok := cpuMap["usage_percent"].(string); ok {
		    val = strings.TrimSuffix(val, "%")
		    if f, err := strconv.ParseFloat(val, 64); err == nil {
		        dashboardData.CPUPercent = f
		    }
		}
	}

	// Memory Usage (string percent -> float)
	var strVal string
	if memoryMap, ok := data["memory"].(map[string]interface{}); ok {
		strVal = memoryMap["used_percent"].(string)
	} else if memoryMap, ok := data["memory"].(map[string]string); ok {
		strVal = memoryMap["used_percent"]
	}
	strVal = strings.TrimSuffix(strVal, "%")
	if fval, err := strconv.ParseFloat(strVal, 64); err == nil {
        dashboardData.MemoryPercent = fval
    }

	// Disk Usage (first disk, first partition)
	if diskArray, ok := data["disk"].([]map[string]interface{}); ok && len(diskArray) > 0 {
		partition:= diskArray[0]["partitions"].([]systemchecks.PartitionInfo)[0]
		dashboardData.DiskRootPercent = partition.UsedPercent
	}

	tmpl, err := template.ParseFiles("webserver/templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl.Execute(w, dashboardData)
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


