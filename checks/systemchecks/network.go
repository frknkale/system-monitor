package systemchecks

import (
	"fmt"
	"monitoring/types"
	"net"
	"net/http"
	"time"
	"os/exec"
	"strconv"
	"strings"
)

func CheckNetwork(cfg types.Config) interface{} {
	if !cfg.Network.Enabled {
		return nil
	}

	result := make(map[string]interface{})

	if cfg.Network.Interfaces {
		ifaces, err := net.Interfaces()
		if err == nil {
			var interfaces []map[string]interface{}
			for _, iface := range ifaces {
				addrs, _ := iface.Addrs()
				var ipList []string
				for _, addr := range addrs {
					ipList = append(ipList, addr.String())
				}

				interfaces = append(interfaces, map[string]interface{}{
					"name":       iface.Name,
					"flags":      iface.Flags.String(),
					"mtu":        iface.MTU,
					"hwaddr":     iface.HardwareAddr.String(),
					"ip_addresses": ipList,
				})
			}
			result["interfaces"] = interfaces
		}
	}

	if len(cfg.Network.Connections) > 0 {
		groupedConns, err := getGroupedConnections(cfg.Network.Connections)
		if err == nil {
			result["active_connections"] = groupedConns
		}
	}

	if len(cfg.Network.ExternalTargets) > 0 {
		result["external_access"] = checkExternalAccess(cfg.Network.ExternalTargets)
	}

	return result
}

func getGroupedConnections(filters []types.ConnectionFilter) ([]map[string]interface{}, error) {
	cmd := "netstat -anp"
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")

	var results []map[string]interface{}

	for _, filter := range filters {
		var filtered []map[string]interface{}

		for _, line := range lines {
			if strings.TrimSpace(line) == "" || !strings.Contains(line, ":") {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 6 {
				continue
			}

			protocol := strings.ToLower(parts[0])
			localAddr := parts[3]
			foreignAddr := parts[4]
			state := parts[5]
			pid := ""
			program := ""

			if len(parts) >= 7 {
				infoParts := strings.Split(parts[6], "/")
				if len(infoParts) == 2 {
					pid = infoParts[0]
					program = infoParts[1]
				}
			}

			localPort := getPortFromAddress(localAddr)

			// Apply filter checks
			if len(filter.Protocols) > 0 && !containsString(filter.Protocols, protocol) {
				continue
			}

			if len(filter.State) > 0 && !containsString(filter.State, state) {
				continue
			}

			if len(filter.Ports) > 0 && !containsIntString(filter.Ports, localPort) {
				continue
			}

			if len(filter.PID) > 0 && !containsIntString(filter.PID, pid) {
				continue
			}

			if len(filter.ProgramNames) > 0 {
				matched := false
				for _, pname := range filter.ProgramNames {
					if strings.Contains(strings.ToLower(program), strings.ToLower(pname)) {
						matched = true
						break
					}
				}
				if !matched {
					continue
				}
			}

			filtered = append(filtered, map[string]interface{}{
				"protocol":        protocol,
				"local_address":   localAddr,
				"foreign_address": foreignAddr,
				"state":           state,
				"pid":             pid,
				"program":         program,
			})
		}

		filterMap := make(map[string]interface{})

		if len(filter.Protocols) > 0 {
			filterMap["Protocols"] = filter.Protocols
		}
		if len(filter.Ports) > 0 {
			filterMap["Ports"] = filter.Ports
		}
		if len(filter.State) > 0 {
			filterMap["State"] = filter.State
		}
		if len(filter.PID) > 0 {
			filterMap["PID"] = filter.PID
		}
		if len(filter.ProgramNames) > 0 {
			filterMap["ProgramNames"] = filter.ProgramNames
		}

		results = append(results, map[string]interface{}{
			"filter":      filterMap,
			"connections": filtered,
		})
	}

	return results, nil
}

func checkExternalAccess(targets []string) []map[string]interface{} {
	var results []map[string]interface{}
	for _, target := range targets {
		entry := map[string]interface{}{
			"target": target,
		}

		var hostname, port = getHostNameAndPort(target)
		
		// DNS Check
		if net.ParseIP(hostname) == nil{
			entry["dns"] = dnsCheck(hostname)
		}

		var httpResult, icmpResult, tcpResult map[string]interface{}


		if strings.Contains(target, "http") && strings.Contains(target, "://"){ 
			httpResult = httpCheck(target)		// HTTP Check
			entry["http"] = httpResult
		} else if port!="" {
			tcpResult = tcpCheck(target)		// TCP Check
			entry["tcp"] = tcpResult
		}

		// ICMP Check
		icmpResult = icmpCheck(hostname)
		entry["icmp"] = icmpResult

		entry["unreachable"] = !(
			(httpResult != nil && httpResult["reachable"] == true) ||
			(tcpResult != nil && tcpResult["reachable"] == true) ||
			(icmpResult != nil && icmpResult["reachable"] == true))

		results = append(results, entry)
	}

	return results
}

func getHostNameAndPort(address string) (string, string) {   // port can be http, domain, etc.
	if strings.Contains(address, "://") {
		address = address[strings.Index(address, "://")+3:]  // Strip the scheme part
	}

	hostname, port, err := net.SplitHostPort(address)
	if err != nil {  	// If no port is specified
		hostname = address
		port = ""
	}

	return hostname, port
}

func dnsCheck(hostname string) map[string]interface{} {
	result := make(map[string]interface{})

	ips, err := net.LookupHost(hostname)
	if err != nil {
		result["dns_resolved"] = false
		result["dns_error"] = err.Error()
	} else {
		result["dns_resolved"] = true
		result["resolved_ips"] = ips
	}

	return result
}

func icmpCheck(hostname string) map[string]interface{} {
	result := make(map[string]interface{})

	_, err := exec.Command("ping", "-c", "2", "-w", "1", hostname).CombinedOutput()
	if err != nil {
		result["reachable"] = false
		result["error"] = fmt.Sprintf("ping error: %s", "unreachable")
	} else {
		result["reachable"] = true
	}

	return result
}


func httpCheck(target string) map[string]interface{} {
	result := make(map[string]interface{})

	resp, err := http.Get(target)
	if err != nil {
		result["reachable"] = false
		result["error"] = err.Error()
	} else {
		result["reachable"] = true
		result["status"] = resp.StatusCode
		resp.Body.Close()
	}

	return  result
}

func tcpCheck(target string) map[string]interface{} {
	result := make(map[string]interface{})

	conn, err := net.DialTimeout("tcp", target, 2 * time.Second)
	if err != nil {
		result["reachable"] = false
		result["error"] = err.Error()
	} else {
		result["reachable"] = true
		conn.Close()
	}

	return result
}

// Extract port from "IP:port"
func getPortFromAddress(address string) string {
	if idx := strings.LastIndex(address, ":"); idx != -1 {
		return address[idx+1:]
	}
	return ""
}

func containsString(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func containsIntString(list []int, item string) bool {
	num, err := strconv.Atoi(item)
	if err != nil {
		return false
	}
	for _, v := range list {
		if v == num {
			return true
		}
	}
	return false
}
