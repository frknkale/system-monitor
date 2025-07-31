package alerter

import (
	"fmt"
	"monitoring/types"
	"time"
	"path/filepath"
	"os"
	"encoding/json"
	"github.com/google/uuid"
)
 
var alertManager *AlertManager

type AlertManager struct {
	Alerts map[string]types.Alert `json:"alerts"`
	OnAlert func(alert types.Alert)
	LogPath string `json:"log_path"`
}
func GetAlertManager() *AlertManager {
	return alertManager
}

func Init(cfg types.Config) {
	if !cfg.Alerter.Enabled {
		alertManager = nil
		return
	}

	alertManager = &AlertManager{
		Alerts: make(map[string]types.Alert),
		LogPath: cfg.Alerter.LogPath,
		OnAlert: func(alert types.Alert) {
			fmt.Printf("New alert: %s at %s\n", alert.Message, alert.FormattedTimestamp)
		},
	}
}

func (alMan *AlertManager) RaiseAlert(message string, status types.HealthStatus, source types.AlerterSources) {
	key := fmt.Sprintf("%s:%s", source, message)
	if existing, ok := alMan.Alerts[key]; ok {
		// fmt.Println("KEY:",key)
		if time.Since(existing.Timestamp) < 2*time.Minute {
			// Alert already exists and is recent, skip raising a new one
			return
		}
	}
	alert := types.Alert{
		ID: 	  uuid.NewString(),
		Timestamp: time.Now(),
		FormattedTimestamp: time.Now().Format("Monday, 02-Jan-06 15:04:05"),
		Message:   message,
		Status:    types.HealthStatus(status),
		Source:    source,
	}
	alMan.Alerts[key] = alert
	if alMan.OnAlert != nil {
		alMan.OnAlert(alert)
	}

	alMan.LogAlerts(alert)

}

func (alMan *AlertManager) LogAlerts(alert types.Alert) {
	if alMan.LogPath == "" {
		return
	}

	_ = os.MkdirAll(filepath.Dir(alMan.LogPath), 0755)

	file, err := os.OpenFile(alMan.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open alert log file: %v\n", err)
		return
	}
	defer file.Close()

	jsonAlert, err := json.Marshal(alert)
	if err != nil {
		fmt.Printf("Failed to marshal alert data: %v\n", err)
		return
	}

	if _, err := file.Write(append(jsonAlert, '\n')); err != nil {
		fmt.Printf("Failed to write alert to log file: %v\n", err)
	}
	fmt.Printf("Alert logged: %s\n", alert.Message)
}
