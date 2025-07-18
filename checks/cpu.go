package checks

import (
	"github.com/shirou/gopsutil/cpu"
	"monitoring/types"
	"time"
)

func CheckCPU(cfg types.Config) interface{} {
	if !cfg.CPU.Enabled {
		return nil
	}

	infoStats, err := cpu.Info()
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	var results []map[string]interface{}
	for _, info := range infoStats {
		result := map[string]interface{}{
			"cpu_number":        info.CPU,
			"vendor_id":  info.VendorID,
			"model_name": info.ModelName,
			"cores":      info.Cores,
			"mhz":        info.Mhz,
			"cache_size": info.CacheSize,
		}
		if len(percentages) > 0 {
			result["usage_percent"] = percentages[0]
		}
		results = append(results, result)
		
	}

	return results
}
