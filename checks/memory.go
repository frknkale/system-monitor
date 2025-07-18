package checks

import (
	"github.com/shirou/gopsutil/mem"
	"monitoring/types"
)

func CheckMemory(cfg types.Config) interface{} {
	for _, memCfg := range cfg.Memory {
		if !memCfg.Enabled {
			continue
		}
		result := make(map[string]interface{})

		if memCfg.Types.RAM {
			vm, err := mem.VirtualMemory()
			if err == nil {
				result["RAM"] = map[string]interface{}{
					"total_mb": vm.Total / 1024 / 1024,
					"used_mb":  vm.Used / 1024 / 1024,
					"free_mb":  vm.Free / 1024 / 1024,
					"used_percent": vm.UsedPercent,
				}
			}
		}
		if memCfg.Types.Swap {
			sm, err := mem.SwapMemory()
			if err == nil {
				result["swap"] = map[string]interface{}{
					"total_mb": sm.Total / 1024 / 1024,
					"used_mb":  sm.Used / 1024 / 1024,
					"free_mb":  sm.Free / 1024 / 1024,
					"used_percent": sm.UsedPercent,
				}
			}
		}
		if memCfg.Types.Total {
			vm, err := mem.VirtualMemory()
			sm, err2 := mem.SwapMemory()
			if err == nil && err2 == nil {
				result["total"] = map[string]interface{}{
					"total_mb": (vm.Total + sm.Total) / 1024 / 1024,
					"used_mb":  (vm.Used + sm.Used) / 1024 / 1024,
					"free_mb":  (vm.Free + sm.Free) / 1024 / 1024,
				}
			}
		}
		return result
	}
	return nil
}
