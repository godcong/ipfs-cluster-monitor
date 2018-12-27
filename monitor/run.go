package monitor

import "context"

type ClusterStatus struct {
	//ipfs info
	//service info
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
