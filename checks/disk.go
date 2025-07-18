package checks

import (
	"sort"
	"github.com/shirou/gopsutil/disk"
	"monitoring/types"
)

func CheckDisk(cfg types.Config) interface{} {
	var results []map[string]interface{}

	for _, diskCfg := range cfg.Disk {
		if !diskCfg.Enabled {
			return nil
		}

		partitions, err := disk.Partitions(true)
		if err != nil {
			return map[string]interface{}{"error": err.Error()}
		}

		filter := diskCfg.Filter

		var filtered []map[string]interface{} 

		for _, p := range partitions {
			usage, err := disk.Usage(p.Mountpoint)
			if err != nil {
				continue
			}
			filtered = append(filtered, map[string]interface{}{
				"device":     p.Device,
				"mountpoint": p.Mountpoint,
				"fstype":     p.Fstype,
				"total_mb":   float64(usage.Total) / 1024 / 1024,
				"used_mb":    float64(usage.Used) / 1024 / 1024,
				"free_mb":	  float64(usage.Free) / 1024 / 1024,
				"used_percent":  float64(usage.UsedPercent),
			})
		}
		sortPartitions := func(items  []map[string]interface{}, key string) {
			keyMap := map[string]string{
				"used_space":       "used_mb",
				"free_space":       "free_mb",
				"total_space":      "total_mb",
				"used_percent":     "used_percent",
			}
			actualKey, ok := keyMap[key]
			if !ok {
				return
			}
			sort.Slice(items, func(i, j int) bool {
				return items[i][actualKey].(float64) > items[j][actualKey].(float64)
			})
		}
		
		if filter.TopDiskSize > 0 {
			sortPartitions(filtered, "total_space")
			filtered = filtered[:filter.TopDiskSize]
		}

		if filter.TopDiskUsage > 0 {
			sortPartitions(filtered, "used_space")
			filtered = filtered[:filter.TopDiskUsage]
		}

		if filter.TopDiskUsagePercent > 0 {
			sortPartitions(filtered, "used_percent")
			filtered = filtered[:filter.TopDiskUsagePercent]
		}

		if filter.TopFreeSpace > 0 {
			sortPartitions(filtered, "free_space")
			filtered = filtered[:filter.TopFreeSpace]
		}
		
		if filter.SortBy != "" {
			sortPartitions(filtered, filter.SortBy)
		}

		results = append(results, map[string]interface{}{
			"filter":    filter,
			"partitions": filtered,
		})
	}
	return results
}
