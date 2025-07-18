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

	percentages, err := cpu.Percent(time.Second, true)
	if err != nil {
		return map[string]interface{}{"error": err.Error()}
	}

	if len(infoStats) == 0 {
		return map[string]interface{}{"error": "no CPU info available"}
	}

	results := map[string]interface{}{
		"model_name": infoStats[0].ModelName,
		"vendor_id":  infoStats[0].VendorID,
		"mhz":        infoStats[0].Mhz,
		"cache_size": infoStats[0].CacheSize,
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

	results["cores"] = cores

	return results
}
