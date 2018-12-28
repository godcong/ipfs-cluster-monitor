package cluster

import "context"

// MonitorStatus ...
type MonitorStatus struct {
}

// StartRun ...
func StartRun(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
			//get info
		}
	}()
}
