package types

import "time"

type Config struct {
	General struct {
		Interval   time.Duration `yaml:"interval"`
		OutputPath string        `yaml:"output_path"`
	} `yaml:"general"`

	Memory []struct {
		Enabled bool `yaml:"enabled"`
		Types   struct {
			RAM   bool `yaml:"RAM"`
			Swap  bool `yaml:"swap"`
			Total bool `yaml:"total"`
		} `yaml:"types"`
	} `yaml:"memory"`

	CPU struct {
		Enabled bool `yaml:"enabled"`
	} `yaml:"cpu"`

	Disk []struct {
		Enabled bool `yaml:"enabled"`
		Filter  struct {
			SortBy				string	`yaml:"sort_by"`
			TopDiskSize 		int		`yaml:"top_disk_size"`
			TopDiskUsage    	int		`yaml:"top_disk_usage"`
			TopDiskUsagePercent	int		`yaml:"top_disk_usage_percent"`
			TopFreeSpace		int		`yaml:"top_free_space"`
		} `yaml:"filter"`
	} `yaml:"disk"`

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
		} `yaml:"filter"`
	} `yaml:"processes"`

	Permissions struct {
		Enabled bool     `yaml:"enabled"`
		Paths   []string `yaml:"paths"`
	} `yaml:"permissions"`

	Network struct {
		Enabled         bool `yaml:"enabled"`
		CheckOpenPorts  bool `yaml:"check_open_ports"`
		CheckInterfaces bool `yaml:"check_interfaces"`
	} `yaml:"network"`
}
