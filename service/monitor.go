package service

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/proto"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Swarm ...
type Swarm struct {
	sync.RWMutex
	address map[string]string
}

// Pin ...
type Pin struct {
	sync.RWMutex
	pins map[string]string
}

// Pins ...
func (p *Pin) Pins() map[string]string {
	p.RLock()
	defer p.RUnlock()
	return p.pins
}

// SetPins ...
func (p *Pin) SetPins(pins map[string]string) {
	p.Lock()
	defer p.Unlock()
	p.pins = pins
}

// Address ...
func (s *Swarm) Address() map[string]string {
	s.RLock()
	defer s.RUnlock()
	return s.address
}

// SetAddress ...
func (s *Swarm) SetAddress(address map[string]string) {
	s.Lock()
	defer s.Unlock()
	s.address = address
}

// Monitor ...
type Monitor struct {
	//isInitialized bool
	//Mode          proto.StartMode
	config        *config.Configure
	context       context.Context
	cancelFunc    context.CancelFunc
	Swarm         *Swarm
	Pin           *Pin
	monitorServer proto.ClusterMonitorClient
}

// Mode ...
func (m *Monitor) Mode() proto.StartMode {
	if m.config != nil {
		return m.config.Monitor.Mode
	}
	return proto.StartMode_Simple
}

// NewMonitor ...
func NewMonitor(cfg *config.Configure) *Monitor {
	return &Monitor{
		//isInitialized: cfg.Initialize,
		//Mode:          cfg.Monitor.Mode,
		config:  cfg,
		context: context.Background(),
		Swarm: &Swarm{
			RWMutex: sync.RWMutex{},
			address: nil,
		},
		Pin: &Pin{
			RWMutex: sync.RWMutex{},
			pins:    nil,
		},
		monitorServer: MonitorClient(NewServerMonitorGRPC(cfg)),
	}
}

// IsInitialized ...
func (m *Monitor) IsInitialized() bool {
	if m.config != nil {
		return m.config.Initialize
	}
	return false
}

// waitingForInitialize ...
func (m *Monitor) waitingForInitialize(ctx context.Context) bool {
	for {
		if !m.IsInitialized() {
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
	log.Printf("client:%+v", *monitor)
	config.SetMonitor(monitor)
	err := cluster.InitMaker(m.config)
	if err == nil {
		return nil
	}
	return xerrors.Errorf("init maker:%w", err)
}

// CustomMaker ...
func (m *Monitor) CustomMaker(monitor *config.Monitor) error {
	log.Printf("client:%+v", *monitor)
	err := cluster.InitMaker(m.config)
	if err == nil {
		return nil
	}

	return xerrors.Errorf("init maker:%w", err)
}

// FileDir ...
func FileDir(path, name string) string {
	return filepath.Join(filepath.Dir(path), name)
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
			if cluster.InitRunning(filepath.Join(m.config.Monitor.Workspace, config.IpfsTmp)) {
				log.Info("init ipfs")
				err := cluster.RunIPFSInit(ctx, m.config)
				if err != nil {
					log.Error(err)
					defer func() { m.Reset() }()
					return
				}
				if m.config.UseCustom {
					log.Info("storage max set with:", m.config.Custom.MaxSize)
					e := cluster.StorageMaxSet(ctx, m.config, m.config.Custom.MaxSize)
					if e != nil {
						log.Error(e)
					}

				}
			}

			if m.Mode() == proto.StartMode_Cluster {
				if cluster.InitRunning(filepath.Join(m.config.Monitor.Workspace, config.ClusterTmp)) {
					log.Println("init ipfs monitor")
					err := cluster.RunServiceInit(ctx, m.config)
					if err != nil {
						log.Error(err)
						defer func() { m.Reset() }()
						return
					}
				}
			}
			if m.Mode() == proto.StartMode_Simple {
				cluster.WriterSwarm(m.config.Monitor.IPFSClient.IpfsPath)
				cluster.RemoveBootstrapIPFS(ctx, m.config)
				cluster.AddBootstrapIPFS(ctx, m.config, m.config.Monitor.IPFSClient.Bootstrap)
			}
			cluster.RunIPFS(ctx, m.config)
			cluster.WaitingIPFS(ctx)

			if m.Mode() == proto.StartMode_Cluster {
				cluster.RunService(ctx, m.config)
				cluster.WaitingService(ctx)
			}

			if m.Mode() == proto.StartMode_Simple {
				m.HandleGRPCAddress(ctx)
				m.Address(ctx)
				m.HandlePins(ctx)
				m.Pins(ctx)
			}

		}
	}()
}

// AddressRes ...
type AddressRes struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail"`
}

// Address ...
func (m *Monitor) Address(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for v := range m.Swarm.Address() {
					cluster.SwarmAddress(v)
				}
				time.Sleep(30 * time.Second)
			}
		}
	}()
}

// HandleGRPCAddress ...
func (m *Monitor) HandleGRPCAddress(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				reply, e := m.monitorServer.MonitorAddress(context.Background(), &proto.MonitorRequest{})
				if e == nil {
					p := make(map[string]string)
					for _, v := range reply.Addresses {
						p[v] = ""
					}
					log.Info(reply.Addresses)
					m.Swarm.SetAddress(p)
				}
				if e != nil {
					log.Error(e)
				}
				time.Sleep(15 * time.Minute)
			}
		}
	}()
}

// HandleAddress ...
func (m *Monitor) HandleAddress(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, e := http.Get("http://cluster.dvbox.net/v0/client/address")
				if e == nil {
					bytes, e := ioutil.ReadAll(resp.Body)
					if e == nil {
						var address AddressRes
						e = jsoniter.Unmarshal(bytes, &address)
						if e == nil {
							log.Info(address.Detail)
							m.Swarm.SetAddress(address.Detail)
						}
					}
				}
				if e != nil {
					log.Error(e)
				}
				time.Sleep(15 * time.Minute)
			}
		}
	}()
}

// PinsRes ...
type PinsRes struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Detail  map[string]string `json:"detail"`
}

// RangePins ...
func RangePins(pins map[string]string) {
	for v := range pins {
		log.Info("pin add:", pins)
		_ = cluster.PinAdd(v)
	}
}

// Pins ...
func (m *Monitor) Pins(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				res, e := cluster.PinLs()
				var pinLs map[string]string
				if e != nil || res.Keys == nil {
					pinLs = m.Pin.Pins()
				} else {
					pinLs = make(map[string]string)
					for v := range m.Pin.Pins() {
						if _, b := (res.Keys)[v]; !b {
							pinLs[v] = "0"
						}
					}
				}
				RangePins(pinLs)
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

// HandleGRPCPins ...
func (m *Monitor) HandleGRPCPins(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				reply, e := m.monitorServer.MonitorPin(context.Background(), &proto.MonitorRequest{})
				if e == nil {
					p := make(map[string]string)
					for _, v := range reply.Pins {
						p[v] = ""
					}
					m.Pin.SetPins(p)
				}

				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

// HandlePins ...
func (m *Monitor) HandlePins(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				resp, e := http.Get("http://cluster.dvbox.net/v0/client/pins")
				if e == nil {
					bytes, e := ioutil.ReadAll(resp.Body)
					if e == nil {
						var pins PinsRes
						e = jsoniter.Unmarshal(bytes, &pins)
						if e == nil {
							log.Info(pins.Detail)
							m.Pin.SetPins(pins.Detail)
						}
					}
				}
				if e != nil {
					log.Error(e)
				}
				time.Sleep(30 * time.Minute)
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
	//err := monitor.RunCMD("rm", env,	"-R", path)
	if err != nil {
		log.Println(err)
	}
	return
}

// Reset ...
func (m *Monitor) Reset() error {
	//reset status
	m.config.Initialize = false

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
