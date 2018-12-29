package cluster

import "context"

// MonitorStatus ...
type MonitorStatus struct {
}

// runMonitor ...
func runMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
		//get info
	}
}
