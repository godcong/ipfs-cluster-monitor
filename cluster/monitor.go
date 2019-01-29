package cluster

import (
	"context"
	"github.com/juju/errors"
	"log"
	"sync"
)

// MonitorPeers ...
const MonitorPeers = "peers"

var monitor sync.Map
var clients sync.Map

// runMonitor ...
func runMonitor(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				peers, err := getPeers()
				if err != nil {
					monitor.Delete(MonitorPeers)

					errors.ErrorStack(err)
					log.Println(err)
					continue
				}
				monitor.Store(MonitorPeers, peers)

			}
			//get info
		}
	}()

}
