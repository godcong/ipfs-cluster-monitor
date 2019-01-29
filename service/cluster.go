package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"log"
	"sync/atomic"
	"time"
)

// ClusterMonitor ...
type ClusterMonitor struct {
	isInitialized bool
	config        *config.Configure
	context       context.Context
	cancelFunc    context.CancelFunc
}

// NewClusterMonitor ...
func NewClusterMonitor(cfg *config.Configure) *ClusterMonitor {
	return &ClusterMonitor{
		config:  cfg,
		context: context.Background(),
	}
}

// Initialized ...
func (m *ClusterMonitor) Initialized() {
	m.isInitialized = true
}

// IsInitialized ...
func (m *ClusterMonitor) IsInitialized() bool {
	return m.isInitialized
}

// waitingForInitialize ...
func (m *ClusterMonitor) waitingForInitialize(ctx context.Context) bool {
	for {
		if !m.isInitialized {
			select {
			case <-ctx.Done():
				return false
			default:
				time.Sleep(m.config.MonitorProperty.Interval)
			}
		}

		return true
	}
}

// Start ...
func (m *ClusterMonitor) Start() {
	var ctx context.Context
	ctx, m.cancelFunc = context.WithCancel(m.context)

	go func() {
		if m.waitingForInitialize(ctx) {

			if m.GetStatus("init") != StatusCreated {
				//wait for make end
				time.Sleep(cfg.MonitorInterval)
			}

			if initCheck(InitIPFS) {
				log.Println("init ipfs")
				m.SetStatus("init", StatusIpfsInit)
				firstRunIPFS()

			}
			if initCheck(InitService) {
				log.Println("init service")
				m.SetStatus("init", StatusServiceInit)
				firstRunService()
			}

			m.SetStatus("init", StatusIpfsRun)
			runIPFS(m.context)
			waitingIpfs(m.context)

			m.SetStatus("init", StatusServiceRun)
			runService(m.context)
			waitingService(m.context)

			if isClient() {
				runJoin(cluster.context)
			} else {
				runMonitor(cluster.context)
			}
			atomic.StoreInt32(&cluster.waiting, -1)
		}
	}()

}

// Stop ...
func (m *ClusterMonitor) Stop() {
	m.stopRunningCMD()
	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}
}
