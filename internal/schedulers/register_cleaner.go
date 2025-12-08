package schedulers

import (
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/api/tools/apm/services_topology"
	"time"
)

func StartClearRegisterScheduler() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				list_services.ClearServicesRegister()
				services_topology.ClearTopoRegister()
			}
		}
	}()
}
