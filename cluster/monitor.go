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
}

// ClientJoin ...
func ClientJoin(key, val string) {
	if _, loaded := clients.LoadOrStore(key, val); loaded {
		log.Println(key, "is already joined")
	} else {
		log.Println(key, "is joined")
	}
}
