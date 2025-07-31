package webserver

import (
	"encoding/json"
	"net/http"
	"html/template"
	"time"
	"fmt"
    "strconv"
    "strings"
	// "reflect"

	"monitoring/types"
	"monitoring/checks"
	"monitoring/cache"
	"monitoring/checks/systemchecks"
	"monitoring/alertcache"
	
	// "bytes"
	// "io/ioutil"
)

var (
	config    types.Config
	sem = make(chan struct{}, 1)
	i = 0
)

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
		fresh := r.URL.Query().Get("fresh") == "true"
		data := getData(fresh)
		// fmt.Println("Available sections:", reflect.ValueOf(data).MapKeys())
		sectionData, ok := data[section]
		if !ok {
			http.Error(w, fmt.Sprintf("Section '%s' not found", section), http.StatusNotFound)
			return
		}

		// Refresh alert cache if needed
		if alertcache.IsExpired() {
			_ = alertcache.Refresh(10)
		}
		alerts := alertcache.Get()

		// Map section name to alert source
		var source types.AlerterSources
		switch section {
		case "cpu":
			source = types.CPU_USAGE_PERCENT
		case "memory":
			source = types.MEM_USAGE_PERCENT
		case "disk":
			source = types.DISK_USAGE_PERCENT
		default:
			source = ""
		}
		filteredAlerts := filterAlerts(alerts, source)

		// Create combined response
		response := map[string]interface{}{
			"metrics": sectionData,
			"alerts":  filteredAlerts,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}


func filterAlerts(alerts []types.Alert, source types.AlerterSources) []types.Alert {
	if len(alerts) == 0 {
		return []types.Alert{}
	}
	var filtered []types.Alert
	for _, a := range alerts {
		if a.Source == source {
			filtered = append(filtered, a)
		}
	}
	return filtered // always non-nil
}

func dismissAlertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}

	err := alertcache.RemoveAlertByID(id)
	if err != nil {
		http.Error(w, "Failed to remove alert: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}



func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	fresh := r.URL.Query().Get("fresh") == "true"
	data := getData(fresh)

	var dashboardData types.DashboardData
	dashboardData.Timestamp = cache.GetLastUpdated().Format("Monday, 02-Jan-06 15:04:05")

	if cpuMap, ok := data["cpu"].(map[string]interface{}); ok {
		if val, ok := cpuMap["usage_percent"].(string); ok {
			val = strings.TrimSuffix(val, "%")
			if f, err := strconv.ParseFloat(val, 64); err == nil {
				dashboardData.CPUPercent = f
			}
		}
	}

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

	if diskArray, ok := data["disk"].([]map[string]interface{}); ok && len(diskArray) > 0 {
		partition := diskArray[0]["partitions"].([]systemchecks.PartitionInfo)[0]
		dashboardData.DiskRootPercent = partition.UsedPercent
	}

	// Load alerts
	if alertcache.IsExpired() {
		_ = alertcache.Refresh(10)
	}
	alerts := alertcache.Get()

	// Render
	tmpl, err := template.ParseFiles("webserver/templates/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	dataStruct := struct {
		types.DashboardData
		Alerts     []types.Alert
	}{
		DashboardData: dashboardData,
		Alerts:        alerts,
	}

	tmpl.Execute(w, dataStruct)
}



func WebServer(cfg types.Config) {
	config = cfg

	alertcache.Init(config.Alerter.LogPath)

	http.HandleFunc("/", dashboardHandler)
	http.HandleFunc("/status", statusHandler)

	http.HandleFunc("/cpu", sectionHandler("cpu"))
	http.HandleFunc("/memory", sectionHandler("memory"))
	http.HandleFunc("/disk", sectionHandler("disk"))
	http.HandleFunc("/network", sectionHandler("network"))
	http.HandleFunc("/services", sectionHandler("services"))
	http.HandleFunc("/permissions", sectionHandler("permissions"))
	http.HandleFunc("/processes", sectionHandler("processes"))
	http.HandleFunc("/alerts/dismiss", dismissAlertHandler)

	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}


