package types

import (
	// "fmt"
	"time"
)

type Config struct {
	General struct {
		Interval   string	`yaml:"interval"`
		OutputPath string 	`yaml:"output_path"`
		LogPath    string	`yaml:"log_path"`
	} `yaml:"general"`

	Memory struct {
		Enabled bool `yaml:"enabled"`
		Total         bool `yaml:"total"`
		Available     bool `yaml:"available"`
		Used          bool `yaml:"used"`
		Free          bool `yaml:"free"`
		UsedPercent   bool `yaml:"used_percent"`
		Active        bool `yaml:"active"`
		Inactive      bool `yaml:"inactive"`
		Buffers       bool `yaml:"buffers"`
		Cached        bool `yaml:"cached"`
		Shared        bool `yaml:"shared"`
		Slab          bool `yaml:"slab"`
		Dirty         bool `yaml:"dirty"`
		SwapTotal     bool `yaml:"swap_total"`
		SwapUsed      bool `yaml:"swap_used"`
		SwapFree      bool `yaml:"swap_free"`
		SwapUsedPercent bool `yaml:"swap_used_percent"`
		SwapIn        bool `yaml:"swap_in"`
		SwapOut       bool `yaml:"swap_out"`
	} `yaml:"memory"`

	CPU struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"cpu"`

	Disk []struct {
		Enabled       bool       `yaml:"enabled"`
		Mounted       bool       `yaml:"mounted"`
		PathsToWatch  []string   `yaml:"paths_to_watch"`
		Filter        struct {
			SortBy             string `yaml:"sort_by"`
			TopDiskSize        int    `yaml:"top_disk_size"`
			TopDiskUsage       int    `yaml:"top_disk_usage"`
			TopDiskUsagePercent int   `yaml:"top_disk_usage_percent"`
			TopFreeSpace       int    `yaml:"top_free_space"`
		}`yaml:"filter"`
	}`yaml:"disk"`

	Processes []struct {
		Enabled bool `yaml:"enabled"`
		Filter  struct {
			RunningHourThreshold int    `yaml:"running_hour_threshold"`
			TopMemoryUsage       int    `yaml:"top_memory_usage"`
			TopCPUUsage          int    `yaml:"top_cpu_usage"`
			State                string `yaml:"state"`
			ParentPID            int32  `yaml:"parent_pid"`
			TTY                  string `yaml:"tty"`
			TopRunningTime       int    `yaml:"top_running_time"`
			LimitProcesses       int    `yaml:"limit_processes"`
		} `yaml:"filter"`
	} `yaml:"processes"`

	Permissions []struct {
		Enabled             bool     `yaml:"enabled"`
		Paths               []string `yaml:"paths"`
		ShowUserPermissions bool     `yaml:"show_user_permissions"`
		CheckUserAccess     []string `yaml:"check_user_access"`
	} `yaml:"permissions"`

	Network struct {
		Enabled     	bool               `yaml:"enabled"`
		Interfaces  	bool               `yaml:"interfaces"`
		Connections 	[]ConnectionFilter `yaml:"connections"`
		ExternalTargets []string       	   `yaml:"external_targets"`
	} `yaml:"network"`

	Alerter struct {
		Enabled bool `yaml:"enabled"`
		LogPath string `yaml:"log_path"`
		AlertSettings struct {
			Memory struct {
				Enabled bool `yaml:"enabled"`
				UsagePercent float64 `yaml:"usage_percent"`
			} `yaml:"memory"`
			CPU struct {
				Enabled	 bool `yaml:"enabled"`
				UsagePercent float64 `yaml:"usage_percent"`
			} `yaml:"cpu"`
			Disk []struct {
				Enabled     bool     `yaml:"enabled"`
				UsagePercent int      `yaml:"usage_percent"`
				PathsToWatch []string `yaml:"paths_to_watch"`
			} `yaml:"disk"`
		} `yaml:"alert_settings"`
	}
}

type ConnectionFilter struct {
    Protocols    []string `yaml:"protocols"`
    Ports        []int    `yaml:"ports"`
    State        []string `yaml:"state"`
    PID          []int    `yaml:"pid"`
    ProgramNames []string `yaml:"program_name"`
}

type HealthStatus string
type AlerterSources string

const (
	HEALTHY   	HealthStatus = "healthy"
	UNHEALTHY 	HealthStatus = "unhealthy"
	ERROR     	HealthStatus = "error"
	
	MEM_USAGE_PERCENT AlerterSources = "memory_usage_percent"
	DISK_USAGE_PERCENT AlerterSources = "disk_usage_percent"
	CPU_USAGE_PERCENT AlerterSources = "cpu_usage_percent"
	CPU_CORE_USAGE_PERCENT AlerterSources = "cpu_core_usage_percent"
)

type Alert struct {
	Timestamp time.Time    `json:"timestamp"`
	Message   string       `json:"message"`
	Status    HealthStatus `json:"status"`
	Source    AlerterSources       `json:"source"`
}


// type AlertManager struct {
// 	Alerts map[string]Alert `json:"alerts"`
// 	OnAlert func(alert Alert)
// 	RaiseAlert func(message string, status HealthStatus, source string)
// }
