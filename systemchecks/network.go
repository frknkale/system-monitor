package systemchecks

import (
	"net"
	"monitoring/types"
)

func CheckNetwork(cfg types.Config) interface{} {
	if !cfg.Network.Enabled {
		return nil
	}

	result := make(map[string]interface{})

	if cfg.Network.CheckInterfaces {
		ifaces, err := net.Interfaces()
		if err == nil {
			var interfaces []map[string]interface{}
			for _, iface := range ifaces {
				interfaces = append(interfaces, map[string]interface{}{
					"name":  iface.Name,
					"flags": iface.Flags.String(),
					"mtu":   iface.MTU,
					"hwaddr": iface.HardwareAddr.String(),
				})
			}
			result["interfaces"] = interfaces
		}
	}

	if cfg.Network.CheckOpenPorts {
		
		ports := []int{}
		
		result["open_ports"] = ports
	}

	return result
}
