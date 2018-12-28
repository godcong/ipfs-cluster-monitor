package cluster

import "context"

// MonitorStatus ...
type MonitorStatus struct {
}

type ipfsInfo struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
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
