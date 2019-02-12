package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"golang.org/x/exp/xerrors"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Monitor ...
type Monitor struct {
	isInitialized bool
	config        *config.Configure
	context       context.Context
	cancelFunc    context.CancelFunc
}

// NewMonitor ...
func NewMonitor(cfg *config.Configure) *Monitor {
	return &Monitor{
		isInitialized: cfg.Initialize,
		config:        cfg,
		context:       context.Background(),
	}
}

// Initialized ...
func (m *Monitor) Initialized() {
	m.isInitialized = true
}

// IsInitialized ...
func (m *Monitor) IsInitialized() bool {

	return m.isInitialized
}

// waitingForInitialize ...
func (m *Monitor) waitingForInitialize(ctx context.Context) bool {
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
func (m *Monitor) InitMaker(monitor *config.Monitor) error {
	log.Printf("monitor:%+v", *monitor)
	m.config.Monitor = *monitor
	err := cluster.InitMaker(m.config)
	if err == nil {
		m.Initialized()
		return nil
	}
	return xerrors.Errorf("init maker:%w", err)
}

// FileDir ...
func FileDir(path, name string) string {
	dir, _ := filepath.Split(path)
	return filepath.Join(dir, name)
}

// Start ...
func (m *Monitor) Start() {
	var ctx context.Context
	ctx, m.cancelFunc = context.WithCancel(m.context)

	go func() {
		if m.waitingForInitialize(ctx) {
			if cluster.InitRunning(m.config.Monitor.IpfsPath) {
				log.Println("init ipfs")
				err := cluster.RunIPFSInit(ctx, m.config)
				if err != nil {
					panic(err)
				}
			}
			if cluster.InitRunning(m.config.Monitor.ClusterPath) {
				log.Println("init ipfs cluster")
				err := cluster.RunServiceInit(ctx, m.config)
				if err != nil {
					panic(err)
				}
			}

			cluster.RunIPFS(ctx, m.config)
			cluster.WaitingIPFS(ctx)

			cluster.RunService(ctx, m.config)
			cluster.WaitingService(ctx)

		}
	}()

}

// Stop ...
func (m *Monitor) Stop() {
	if m.cancelFunc != nil {
		m.cancelFunc()
		m.cancelFunc = nil
	}
}

// Clear ...
func clear(path string) {
	log.Println("clear", path)
	err := os.RemoveAll(path)
	//err := cluster.RunCMD("rm", env,	"-R", path)
	if err != nil {
		log.Println(err)
	}
	return
}

// Reset ...
func (m *Monitor) Reset() error {
	//reset status
	m.isInitialized = false

	//stop running ipfs and service
	m.Stop()
	time.Sleep(15 * time.Second)
	log.Println("stop all")

	clear(config.IpfsPath())
	clear(config.IpfsClusterPath())
	err := os.Remove(m.config.FD())
	log.Println(err)
	//dir, name := m.config.Root, m.config.ConfigName

	//reset config
	//m.config.Initialize = false
	//m.config = config.DefaultConfig()
	//m.config.Root, m.config.ConfigName = dir, name

	err = m.InitMaker(&m.config.Monitor)
	log.Println(err)

	//rerun
	go m.Start()
	log.Println("starting")
	time.Sleep(15 * time.Second)
	m.isInitialized = true

	return nil
}
