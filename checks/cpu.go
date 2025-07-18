package checks

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/load"
	"monitoring/types"
	"time"
	"fmt"
)

func CheckCPU(cfg types.Config) interface{} {
	if !cfg.CPU.Enabled {
		return nil
	}

	infoStats, err := cpu.Info()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}
	result := map[string]interface{}{
		"model_name": infoStats[0].ModelName,
		"vendor_id":  infoStats[0].VendorID,
		"mhz":        infoStats[0].Mhz,
		"cache_size": infoStats[0].CacheSize,
	}

	percentages, err := cpu.Percent(time.Second, true)
    if err != nil {
        result["error"] = fmt.Sprintf("Failed to get CPU percentages: %v", err)
        return result
    }
	
    loads, err := load.Avg()
    if err != nil {
        result["error"] = fmt.Sprintf("Failed to get load averages: %v", err)
        return result
    } 

    times, err := cpu.Times(false) 
    if err != nil {
        result["error"] = fmt.Sprintf("Failed to get CPU times: %v", err)
        return result
    }
    if len(times) == 0 {
        result["error"] = "CPU times returned empty slice"
        return result
    }

	var cores []map[string]interface{}
	for i, info := range infoStats {
		coreData := map[string]interface{}{
			"cpu_number": info.CPU,
		}
		if len(percentages) > i {
			coreData["usage_percent"] = percentages[i]
		}
		cores = append(cores, coreData)
	}

	result["cores"] = cores

	result["time_spent_user"] = fmt.Sprintf("%.2f", times[0].User)
    result["time_spent_system"] = fmt.Sprintf("%.2f", times[0].System)
    result["idle"] = fmt.Sprintf("%.2f", times[0].Idle)
    // result["iowait"] = fmt.Sprintf("%.2f", times[0].Iowait)
    // result["irq"] = fmt.Sprintf("%.2f", times[0].Irq)
    // result["softirq"] = fmt.Sprintf("%.2f", times[0].Softirq)
    // result["steal"] = fmt.Sprintf("%.2f", times[0].Steal)
    // result["guest"] = fmt.Sprintf("%.2f", times[0].Guest)
    // result["guest_nice"] = fmt.Sprintf("%.2f", times[0].GuestNice)
	result["usage_percent"] = fmt.Sprintf("%.2f", percentages[0])
	result["load_1min"] = fmt.Sprintf("%.2f", loads.Load1)
	result["load_5min"] = fmt.Sprintf("%.2f", loads.Load5)
	result["load_15min"] = fmt.Sprintf("%.2f", loads.Load15)

	return result
}
