package alerter

import (
	"fmt"
	"monitoring/types"
)

func EvaluateAlerts(results []types.CheckResult) {
    for _, r := range results {
		if r.Status == "error" {
			handleError(r)
		} else if r.Status == "unhealthy" {
			handleUnhealthy(r)
		} else if r.Status == "healthy" {
			handleHealthy(r)
		}
}

