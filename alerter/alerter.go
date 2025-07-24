package alerter

import (
	"fmt"
	"monitoring/types"
	"time"
	"path/filepath"
	"os"
	"encoding/json"
)

type AlertManager struct {
	Alerts map[string]types.Alert `json:"alerts"`
	OnAlert func(alert types.Alert)
	LogPath string `json:"log_path"`
	// RaiseAlert func(message string, status types.HealthStatus, source string)
}

func NewAlertManager(logPath string) *AlertManager {
	return &AlertManager{
		Alerts: make(map[string]types.Alert),
		LogPath: logPath,
		OnAlert: func(alert types.Alert) {
			fmt.Printf("New alert: %s at %s\n", alert.Message, alert.Timestamp)
		},
	}
}

func (alMan *AlertManager) RaiseAlert(message string, status types.HealthStatus, source types.AlerterSources) {
	key := fmt.Sprintf("%s:%s", source, message)
	if existing, ok := alMan.Alerts[key]; ok {
		if time.Since(existing.Timestamp) < 3*time.Minute {
			// Alert already exists and is recent, skip raising a new one
			return
		}
	}
	alert := types.Alert{
		Timestamp: time.Now(),
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

func AlerterHandler(cfg types.Config) *AlertManager {
	if !cfg.Alerter.Enabled {
		return nil
	}
	return NewAlertManager(cfg.Alerter.LogPath)
}
