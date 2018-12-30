package cluster

import (
	"context"
	"github.com/juju/errors"
	"log"
	"sync"
	"time"
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
					time.Sleep(cfg.MonitorInterval)
					errors.ErrorStack(err)
					log.Println(err)
					continue
				}
				monitor.Store(MonitorPeers, peers)
				time.Sleep(cfg.MonitorInterval)
			}
			//get info
		}
	}()

}
