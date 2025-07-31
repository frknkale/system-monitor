package systemchecks

import (
	"fmt"
	"sort"
	"github.com/shirou/gopsutil/v3/disk"
	"monitoring/types"
	"monitoring/alerter"
)

type PartitionInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Fstype      string  `json:"fstype"`
	TotalMB     float64 `json:"total_mb"`
	UsedMB      float64 `json:"used_mb"`
	FreeMB      float64 `json:"free_mb"`
	UsedPercent float64 `json:"used_percent"`
	Error       string  `json:"error,omitempty"`
}



func CheckDisk(config types.Config) interface{} {	
	var results []map[string]interface{}

	for _, cfg := range config.Disk {
		if !cfg.Enabled {
			continue
		}
		filter := cfg.Filter

		var paths []string
		partMap := make(map[string]string)

		allPartitions, _ := disk.Partitions(!cfg.Mounted)
		for _, part := range allPartitions {
			partMap[part.Mountpoint] = part.Device
		}

		if len(cfg.PathsToWatch) == 1 && cfg.PathsToWatch[0] == "*" {
			for _, part := range allPartitions {
				paths = append(paths, part.Mountpoint)
			}
		} else {
			paths = cfg.PathsToWatch
		}

		var partitions []PartitionInfo
		for _, mount := range paths {
			usage, err := disk.Usage(mount)
			device := partMap[mount]
			if err != nil {
				partitions = append(partitions, PartitionInfo{
					Device:     device,
					Mountpoint: mount,
					Error:      err.Error(),
				})
				continue
			}

			partitions = append(partitions, PartitionInfo{
				Device:      device,
				Mountpoint:  usage.Path,
				Fstype:      usage.Fstype,
				TotalMB:     float64(usage.Total) / 1024 / 1024,
				UsedMB:      float64(usage.Used) / 1024 / 1024,
				FreeMB:      float64(usage.Free) / 1024 / 1024,
				UsedPercent: usage.UsedPercent,
			})

			for _, alertRule := range config.Alerter.AlertSettings.Disk {
				if !alertRule.Enabled {
					continue
				}
			
				watchSet := make(map[string]bool)
				for _, p := range alertRule.PathsToWatch {
					watchSet[p] = true
				}
				
				// usage.UsedPercent = 81
				
				alerter := alerter.GetAlertManager()
				
				if alerter!=nil && usage.UsedPercent > float64(alertRule.UsagePercent) && watchSet[usage.Path] {
					alerter.RaiseAlert(
						fmt.Sprintf("Disk usage is above the given threshold: %.2f%% used on %s", usage.UsedPercent, usage.Path),
						types.UNHEALTHY,
						types.DISK_USAGE_PERCENT,
					)
				}
			}
			
			// Alert if disk usage percent exceeds threshold
			// if alerter != nil &&
			// 	 int(usage.UsedPercent) > config.Alerter.AlertThresholds.Disk.UsagePercent && 
			// 	 {
			// 	alerter.RaiseAlert(
			// 		fmt.Sprintf("Disk usage is above the given threshold: %.2f%% used on %s", usage.UsedPercent, usage.Path),
			// 		types.UNHEALTHY,
			// 		"disk_used_percent",
			// 	)
			// }
		}

		sortPartitions := func(key string) {
			keyMap := map[string]func(p PartitionInfo) float64{
				"used_space":   func(p PartitionInfo) float64 { return p.UsedMB },
				"free_space":   func(p PartitionInfo) float64 { return p.FreeMB },
				"total_space":  func(p PartitionInfo) float64 { return p.TotalMB },
				"used_percent": func(p PartitionInfo) float64 { return p.UsedPercent },
			}

			fn, ok := keyMap[key]
			if !ok {
				return
			}

			sort.Slice(partitions, func(i, j int) bool {
				return fn(partitions[i]) > fn(partitions[j])
			})
		}
		
		filterMap := map[string]interface{}{}

		if filter.TopDiskSize > 0 {
			filterMap["TopDiskSize"] = filter.TopDiskSize
			sortPartitions("total_space")
			partitions = partitions[:min(filter.TopDiskSize, len(partitions))]
		}
		if filter.TopDiskUsage > 0 {
			filterMap["TopDiskUsage"] = filter.TopDiskUsage
			sortPartitions("used_space")
			partitions = partitions[:min(filter.TopDiskUsage, len(partitions))]
		}
		if filter.TopDiskUsagePercent > 0 {
			filterMap["TopDiskUsagePercent"] = filter.TopDiskUsagePercent
			sortPartitions("used_percent")
			partitions = partitions[:min(filter.TopDiskUsagePercent, len(partitions))]
		}
		if filter.TopFreeSpace > 0 {
			filterMap["TopFreeSpace"] = filter.TopFreeSpace
			sortPartitions("free_space")
			partitions = partitions[:min(filter.TopFreeSpace, len(partitions))]
		}
		if filter.SortBy != "" {
			filterMap["SortBy"] = filter.SortBy
			sortPartitions(filter.SortBy)
		}

		results = append(results, map[string]interface{}{
			"filter":    filterMap,
			"partitions": partitions,
			"paths_to_watch": cfg.PathsToWatch,
			"mounted":        cfg.Mounted,
		})
	}

	return results

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
