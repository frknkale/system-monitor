package alertcache

import (
	"bufio"
	"encoding/json"
	"os"
	"sync"
	"time"
	"fmt"

	"monitoring/types"
)

var (
	alerts      []types.Alert
	alertMutex  sync.RWMutex
	lastUpdated time.Time
	cacheDuration = 20 * time.Second
	logPath     string
)

// Initialize the cache with the alert file path
func Init(path string) {
	logPath = path
}

// Check if the cache is expired
func IsExpired() bool {
	alertMutex.RLock()
	defer alertMutex.RUnlock()
	return time.Since(lastUpdated) > cacheDuration
}

// Get a copy of cached alerts
func Get() []types.Alert {
	alertMutex.RLock()
	defer alertMutex.RUnlock()

	// shallow copy
	copyAlerts := make([]types.Alert, len(alerts))
	copy(copyAlerts, alerts)
	return copyAlerts
}

func RemoveAlertByID(id string) error {
	alertMutex.Lock()
	defer alertMutex.Unlock()

	var updated []types.Alert
	var found bool

	for _, alert := range alerts {
		if alert.ID != id {
			updated = append(updated, alert)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("alert with ID %s not found", id)
	}

	// Rewrite log file
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, a := range updated {
		b, err := json.Marshal(a)
		if err != nil {
			continue
		}
		writer.Write(b)
		writer.WriteByte('\n')
	}
	writer.Flush()

	alerts = updated
	lastUpdated = time.Now()
	return nil
}


func RemoveAlertByTimestamp(ts string) error {
	alertMutex.Lock()
	defer alertMutex.Unlock()

	var updated []types.Alert
	var removed bool

	for _, alert := range alerts {
		t := alert.Timestamp.Format(time.RFC3339)
		if t != ts {
			updated = append(updated, alert)
		} else {
			removed = true
		}
	}

	if !removed {
		fmt.Println("Alert not found for timestamp:", ts)
		return nil
	}

	// Write updated list back to file
	file, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, a := range updated {
		b, err := json.Marshal(a)
		if err != nil {
			continue
		}
		writer.Write(b)
		writer.WriteByte('\n')
	}
	writer.Flush()

	alerts = updated
	lastUpdated = time.Now()

	fmt.Println("Removed alert with timestamp:", ts)
	return nil
}



// Update the cache from file
func Refresh(maxLines int) error {
	file, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	var newAlerts []types.Alert
	for _, line := range lines {
		var a types.Alert
		if err := json.Unmarshal([]byte(line), &a); err == nil {
			newAlerts = append(newAlerts, a)
		}
	}

	alertMutex.Lock()
	alerts = newAlerts
	lastUpdated = time.Now()
	alertMutex.Unlock()

	return nil
}
