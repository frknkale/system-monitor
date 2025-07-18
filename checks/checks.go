package checks

import (
	"monitoring/types"
)

func RunAllChecks(cfg types.Config) map[string]interface{} {
	result := make(map[string]interface{})

	if mem := CheckMemory(cfg); mem != nil {
		result["memory"] = mem
	}
	if cpu := CheckCPU(cfg); cpu != nil {
		result["cpu"] = cpu
	}
	if disk := CheckDisk(cfg); disk != nil {
		result["disk"] = disk
	}
	if procs := CheckProcesses(cfg); procs != nil {
		result["processes"] = procs
	}
	if perms := CheckPermissions(cfg); perms != nil {
		result["permissions"] = perms
	}
	if netw := CheckNetwork(cfg); netw != nil {
		result["network"] = netw
	}

	return result
}
