package systemchecks

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/mem"
	"monitoring/types"
	"monitoring/alerter"
)


func CheckMemory(cfg types.Config) map[string]string {
	alerter := alerter.AlerterHandler(cfg)

	if !cfg.Memory.Enabled {
		return nil
	}

	config := cfg.Memory
	result := make(map[string]string)

	format := func(b uint64) string {
		gb := float64(b) / (1024 * 1024 * 1024)
		kb := b / 1024
		return fmt.Sprintf("%.2f GB (%d KB)", gb, kb)
	}

	// formatGB := func(b uint64) float64 {
	// 	return float64(b) / (1024 * 1024 * 1024)
	// }

	// Physical memory
	vm, err := mem.VirtualMemory()
	if err != nil {
		result["error"] = fmt.Sprintf("VirtualMemory error: %v", err)
		return result
	}

	if config.Total {
		result["total"] = format(vm.Total)
	}
	if config.Available {
		result["available"] = format(vm.Available)
	}
	if config.Used {
		result["used"] = format(vm.Used)
	}
	if config.Free {
		result["free"] = format(vm.Free)
	}
	if config.UsedPercent {
		result["used_percent"] = fmt.Sprintf("%.2f%%", vm.UsedPercent)

		// dummy := 83

		if alerter != nil && vm.UsedPercent > cfg.Alerter.AlertThresholds.Memory.UsagePercent {
			alerter.RaiseAlert(
				fmt.Sprintf("Memory usage is above the given threshold: %.2f%% used", vm.UsedPercent),
				types.UNHEALTHY,
				types.MEM_USAGE_PERCENT,
			)
		}
	}

	if config.Active {
		result["active"] = format(vm.Active)
	}
	if config.Inactive {
		result["inactive"] = format(vm.Inactive)
	}
	if config.Buffers {
		result["buffers"] = format(vm.Buffers)
	}
	if config.Cached {
		result["cached"] = format(vm.Cached)
	}
	if config.Shared {
		result["shared"] = format(vm.Shared)
	}
	if config.Slab {
		result["slab"] = format(vm.Slab)
	}
	if config.Dirty {
		result["dirty"] = format(vm.Dirty)
	}

	// Swap memory
	sm, err := mem.SwapMemory()
	if err == nil {
		if config.SwapTotal {
			result["swap_total"] = format(sm.Total)
		}
		if config.SwapUsed {
			result["swap_used"] = format(sm.Used)
		}
		if config.SwapFree {
			result["swap_free"] = format(sm.Free)
		}
		if config.SwapUsedPercent {
			result["swap_used_percent"] = fmt.Sprintf("%.2f%%", sm.UsedPercent)
		}
		if config.SwapIn {
			result["swap_in"] = fmt.Sprintf("%d pages", sm.Sin)
		}
		if config.SwapOut {
			result["swap_out"] = fmt.Sprintf("%d pages", sm.Sout)
		}
	} else {
		result["swap_error"] = fmt.Sprintf("SwapMemory error: %v", err)
	}

	return result
}

