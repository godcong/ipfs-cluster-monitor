package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"time"
)

// Monitor ...
type Monitor struct {
	isInitialized bool
	Mode          proto.StartMode
	config        *config.Configure
	context       context.Context
	cancelFunc    context.CancelFunc
}

// NewMonitor ...
func NewMonitor(cfg *config.Configure) *Monitor {
	return &Monitor{
		isInitialized: cfg.Initialize,
		Mode:          cfg.Monitor.Mode,
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
	if !m.config.Monitor.Enable {
		return
	}
	var ctx context.Context
	ctx, m.cancelFunc = context.WithCancel(m.context)

	go func() {

		if m.waitingForInitialize(ctx) {
			if cluster.InitRunning(filepath.Join(m.config.Monitor.Workspace, config.Ipfs)) {
				log.Println("init ipfs")
				err := cluster.RunIPFSInit(ctx, m.config)
				if err != nil {
					log.Error(err)
					defer func() { m.Reset() }()
					return
				}
			}

			if m.Mode == proto.StartMode_Cluster {
				if cluster.InitRunning(filepath.Join(m.config.Monitor.Workspace, config.Cluster)) {
					log.Println("init ipfs cluster")
					err := cluster.RunServiceInit(ctx, m.config)
					if err != nil {
						log.Error(err)
						defer func() { m.Reset() }()
						return
					}
				}
			}

			cluster.RunIPFS(ctx, m.config)
			cluster.WaitingIPFS(ctx)

			if m.Mode == proto.StartMode_Cluster {
				cluster.RunService(ctx, m.config)
				cluster.WaitingService(ctx)
			}

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
	clear(config.ClusterPath())
	err := os.Remove(m.config.FD())
	log.Println(err)
	//dir, name := m.config.ConfigPath, m.config.ConfigName

	//reset config
	//m.config.Initialize = false
	//m.config = config.DefaultConfig()
	//m.config.ConfigPath, m.config.ConfigName = dir, name

	err = m.InitMaker(&m.config.Monitor)
	log.Println(err)

	//rerun
	m.Start()
	log.Println("starting")
	time.Sleep(15 * time.Second)
	//m.isInitialized = true

	return nil
}
