package checks

import (
	"monitoring/types"
	"monitoring/systemchecks"
	// "monitoring/servicechecks"
)

func RunAllChecks(cfg types.Config) map[string]interface{} {
	result := make(map[string]interface{})

	if mem := systemchecks.CheckMemory(cfg); mem != nil {
		result["memory"] = mem
	}
	if cpu := systemchecks.CheckCPU(cfg); cpu != nil {
		result["cpu"] = cpu
	}
	if disk := systemchecks.CheckDisk(cfg); disk != nil {
		result["disk"] = disk
	}
	if procs := systemchecks.CheckProcesses(cfg); procs != nil {
		result["processes"] = procs
	}
	if perms := systemchecks.CheckPermissions(cfg); perms != nil {
		result["permissions"] = perms
	}
	if netw := systemchecks.CheckNetwork(cfg); netw != nil {
		result["network"] = netw
	}

	return result
}
