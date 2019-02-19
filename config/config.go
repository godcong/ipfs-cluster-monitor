package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

// DefaultFileName ...
const DefaultFileName = "monitor.json"

// Ipfs ...
const Ipfs = "ipfs"

// Cluster ...
const Cluster = "cluster"

// Database ...
type Database struct {
	Prefix   string `toml:"prefix"`
	Type     string `toml:"type"`
	Addr     string `toml:"addr"`
	Port     string `toml:"port"`
	Password string `toml:"password"`
	Username string `toml:"username"`
	DB       string `toml:"db"`
}

// Callback ...
type Callback struct {
	Type     string `toml:"type"`
	BackType string `toml:"back_type"`
	BackAddr string `toml:"back_addr"`
}

// IPFS ...
type IPFS struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

// GRPC ...
type GRPC struct {
	Enable bool   `toml:"enable"`
	Type   string `toml:"type"`
	Path   string `toml:"path"`
	Port   string `toml:"port"`
}

// REST ...
type REST struct {
	Enable  bool   `toml:"enable"`
	Type    string `toml:"type"`
	Path    string `toml:"path"`
	BackURL string `toml:"back_url"`
	Port    string `toml:"port"`
}

// Queue ...
type Queue struct {
	Type     string `toml:"type"`
	HostPort string `toml:"host_port"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

// HostInfo ...
type HostInfo struct {
	Type    string `toml:"type"`
	Addr    string `toml:"addr"`
	Port    string `toml:"port"`
	Version string `toml:"version"`
}

// Requester ...
type Requester struct {
	Type string `toml:"type"`
}

// Monitor ...
type Monitor struct {
	Token       string   `toml:"token"`
	Enable      bool     `toml:"enable"`
	Type        string   `toml:"type"`
	Addr        string   `toml:"addr"`
	Port        string   `toml:"port"`
	Workspace   string   `toml:"workspace"`
	Secret      string   `toml:"secret"`
	Bootstrap   []string `toml:"bootstrap"`
	IpfsPath    string   `toml:"ipfs_path"`
	ClusterPath string   `toml:"cluster_path"`
}

// MustMonitor ...
func MustMonitor(secret, boot, workspace string) *Monitor {
	ipfs := ""
	cluster := ""
	if workspace != "" {
		log.Println(workspace)
		ipfs = filepath.Join(workspace, "data", Ipfs)
		cluster = filepath.Join(workspace, "data", Cluster)
	}

	return &Monitor{
		Enable:      true,
		Workspace:   workspace,
		Secret:      DefaultString(secret, "27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34"),
		Bootstrap:   []string{DefaultString(boot, "/ip4/47.101.169.94/tcp/9096/ipfs/QmeQzPKd7HzKZwBKNmnJnyub3YyCBvtcWraaJKEKk1BWmx")},
		IpfsPath:    DefaultString(ipfs, HomePath(".ipfs")),
		ClusterPath: DefaultString(cluster, HomePath(".ipfs-cluster")),
	}
}

// Env ...
func (m *Monitor) Env() (env []string) {
	var e error
	e = os.Setenv("IPFS_PATH", string(m.IpfsPath))
	if e != nil {
		panic(e)
	}

	e = os.Setenv("CLUSTER_SECRET", string(m.Secret))
	if e != nil {
		panic(e)
	}

	e = os.Setenv("IPFS_CLUSTER_PATH", string(m.ClusterPath))
	if e != nil {
		panic(e)
	}

	env = os.Environ()
	log.Println(env)
	return
}

// MonitorProperty ...
type MonitorProperty struct {
	Version             string        `toml:"version"`
	IpfsCommandName     string        `toml:"ipfs_command_name"`
	ClusterCommandName  string        `toml:"cluster_command_name"`
	RemoteAddrPort      string        `toml:"remote_addr_port"`
	Interval            time.Duration `toml:"interval"`
	ServerCheckInterval time.Duration `toml:"server_check_interval"`
	MonitorInterval     time.Duration `toml:"monitor_interval"`
	ResetWaiting        int           `toml:"reset_waiting"`
}

// Configure ...
type Configure struct {
	Mode            int             `toml:"mode"`
	Initialize      bool            `toml:"-"`
	Root            string          `toml:"-"`
	ConfigName      string          `toml:"-"`
	Monitor         Monitor         `toml:"monitor"`
	MonitorProperty MonitorProperty `toml:"monitor_property"`
	GRPC            GRPC            `toml:"grpc"`
	REST            REST            `toml:"rest"`
	IPFS            IPFS            `toml:"ipfs"`
	Requester       Requester       `toml:"requester"`
	Callback        Callback        `toml:"callback"`
}

var config *Configure

// Initialize ...
func Initialize(filePath ...string) error {
	//if filePath == nil {
	//	filePath = []string{"config.toml"}
	//}
	log.Println(filePath)
	config = LoadConfig(filePath[0])
	config.Root, config.ConfigName = filepath.Split(filePath[0])
	return nil
}

// IsExists ...
func IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// LoadConfig ...
func LoadConfig(filePath string) *Configure {
	cfg := DefaultConfig()

	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		log.Println("open:", err)
		return cfg
	}
	defer openFile.Close()
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(cfg)
	if err != nil {
		log.Println("decode:", err)
		return cfg
	}
	cfg.Initialize = true
	log.Printf("config: %+v", cfg)
	return cfg
}

// HomePath ...
func HomePath(name string) string {
	// We try guessing user's home from the HOME variable. This
	// allows HOME hacks for things like Snapcraft builds. HOME
	// should be set in all UNIX by the OS. Alternatively, we fall back to
	// usr.HomeDir (which should work on Windows etc.).
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("cannot get current user: %s", err))
		}
		home = usr.HomeDir
	}

	return filepath.Join(home, name)
}

// IpfsPath ...
func IpfsPath() string {
	if config.Monitor.IpfsPath != "" {
		return config.Monitor.IpfsPath
	}
	return HomePath(".ipfs")
}

// IpfsClusterPath ...
func IpfsClusterPath() string {
	if config.Monitor.ClusterPath != "" {
		return config.Monitor.ClusterPath
	}
	return HomePath(".ipfs-cluster")
}

// DefaultConfig ...
func DefaultConfig() *Configure {
	return &Configure{
		Initialize: false,
		Root:       "",
		Monitor: Monitor{
			Enable:    true,
			Type:      "tcp",
			Addr:      "localhost",
			Port:      ":7784",
			Token:     "2UQEoTCGV7j689CEImQDAjcv7k1X0ZpxT2yzCX8vaqRg1vKp5f0uScvPVB7yuZPP",
			Secret:    "27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34",
			Bootstrap: []string{"/ip4/47.101.169.94/tcp/9096/ipfs/QmeQzPKd7HzKZwBKNmnJnyub3YyCBvtcWraaJKEKk1BWmx"},
			//Workspace: "",
			//IpfsPath:    HomePath(".ipfs"),
			//ClusterPath: HomePath(".ipfs-cluster"),
		},
		MonitorProperty: MonitorProperty{
			Version:             "",
			IpfsCommandName:     "/data/local/bin/ipfs",
			ClusterCommandName:  "/data/local/bin/ipfs-cluster-service",
			RemoteAddrPort:      "",
			Interval:            3 * time.Second,
			ServerCheckInterval: 3 * time.Second,
			MonitorInterval:     3 * time.Second,
			ResetWaiting:        0,
		},
		GRPC: GRPC{
			Enable: true,
			Type:   "",
			Path:   "",
			Port:   "",
		},
		IPFS: IPFS{
			Host: "",
			Port: "",
		},
	}
}

// Config ...
func Config() *Configure {
	if config == nil {
		panic("nil config")
	}
	return config
}

// FD config file dictionary
func (c *Configure) FD() string {
	return filepath.Join(c.Root, c.ConfigName)
}

// DefaultString ...
func DefaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// CheckExist ...
func (c *Configure) CheckExist() bool {
	_, err := os.Stat(filepath.Join(c.Root, DefaultFileName))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// DefaultMonitorProperty ...
func DefaultMonitorProperty() *MonitorProperty {
	return &MonitorProperty{
		Version:             "v0",
		IpfsCommandName:     "ipfs",
		ClusterCommandName:  "ipfs-cluster-service",
		RemoteAddrPort:      "127.0.0.1:7758",
		Interval:            1 * time.Second,
		MonitorInterval:     5 * time.Second,
		ServerCheckInterval: 60 * time.Second,
		ResetWaiting:        30,
	}
}

// CheckError ...
func CheckError(err error) error {
	return xerrors.Errorf("check err:%w", err)
}
