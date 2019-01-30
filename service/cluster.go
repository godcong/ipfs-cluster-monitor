package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"golang.org/x/exp/xerrors"
	"log"
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
		isInitialized: false,
		config:        cfg,
		context:       context.Background(),
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
				log.Println("waiting for init")
				continue
			}
		}

		return true
	}
}

// InitMaker ...
func (m *ClusterMonitor) InitMaker(monitor *config.Monitor) error {
	err := cluster.InitMaker(m.config, m.config.Root)
	if err == nil {
		m.Initialized()
		return nil
	}
	return xerrors.Errorf("init maker:%w", err)
}

// Start ...
func (m *ClusterMonitor) Start() {
	var ctx context.Context
	ctx, m.cancelFunc = context.WithCancel(m.context)

	go func() {
		if m.waitingForInitialize(ctx) {

			if cluster.InitRunning(config.IpfsPath()) {
				log.Println("init ipfs")
				err := cluster.RunIPFSInit(ctx, m.config)
				if err != nil {
					panic(err)
				}
			}
			if cluster.InitRunning(config.IpfsClusterPath()) {
				log.Println("init ipfs cluster")
				err := cluster.RunServiceInit(ctx, m.config)
				if err != nil {
					panic(err)
				}
			}

			cluster.RunIPFS(m.context, m.config.Monitor.Env())
			cluster.WaitingIPFS(m.context)

			cluster.RunService(m.context, m.config)
			cluster.WaitingService(m.context)

		}
	}()

}

// Stop ...
func (m *ClusterMonitor) Stop() {
	//m.stopRunningCMD()
	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}
}
