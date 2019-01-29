package config

import (
	"fmt"
	"github.com/juju/errors"
	"github.com/pelletier/go-toml"
	"golang.org/x/exp/xerrors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

// DefaultFileName ...
const DefaultFileName = "monitor.json"

// InitIPFS ...
const InitIPFS = "ipfs"

// InitCluster ...
const InitCluster = "cluster"

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

// Media ...
type Media struct {
	Upload      string `toml:"upload"`        //上传路径
	Transfer    string `toml:"transfer"`      //转换路径
	M3U8        string `toml:"m3u8"`          //m3u8文件名
	KeyURL      string `toml:"key_url"`       //default url
	KeyDest     string `toml:"key_dest"`      //key 文件输出目录
	KeyFile     string `toml:"key_file"`      //key文件名
	KeyInfoFile string `toml:"key_info_file"` //keyFile文件名
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
	Secret      string `toml:"secret"`
	Bootstrap   string `toml:"bootstrap"`
	Path        string `toml:"path"`
	ClusterPath string `toml:"cluster_path"`
}

// Env ...
func (m *Monitor) Env() (env []string) {
	env = append(env, strings.Join([]string{"IPFS_PATH", string(m.Path)}, "="))
	env = append(env, strings.Join([]string{"CLUSTER_SECRET", string(m.Secret)}, "="))
	env = append(env, strings.Join([]string{"IPFS_CLUSTER_PATH", string(m.ClusterPath)}, "="))
	return
}

// HostType ...
type HostType string

// ClassHost ...
var (
	HostServer HostType = "server"
	HostClient HostType = "client"
)

// MonitorProperty ...
type MonitorProperty struct {
	Version             string
	RootPath            string
	CommandName         string
	ServiceCommandName  string
	HostType            HostType
	RemoteIP            string
	RemotePort          string
	Interval            time.Duration
	ServerCheckInterval time.Duration
	MonitorInterval     time.Duration
	ResetWaiting        int
}

// Configure ...
type Configure struct {
	Root            string          `toml:"root"`
	Monitor         Monitor         `toml:"monitor"`
	MonitorProperty MonitorProperty `toml:"monitor_property"`
	Database        Database        `toml:"database"`
	Censor          HostInfo        `toml:"censor"`
	Node            HostInfo        `toml:"node"`
	Media           Media           `toml:"media"`
	Queue           Queue           `toml:"queue"`
	GRPC            GRPC            `toml:"grpc"`
	REST            REST            `toml:"rest"`
	IPFS            IPFS            `toml:"ipfs"`

	Requester Requester `toml:"requester"`
	Callback  Callback  `toml:"callback"`
}

var config *Configure

// Initialize ...
func Initialize(filePath ...string) error {
	if filePath == nil {
		filePath = []string{"config.toml"}
	}
	config = LoadConfig(filePath[0])
	config.Root = filePath[0]
	return nil
}

// IsExists ...
func IsExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Panicln(err)
	}
	return true
}

// LoadConfig ...
func LoadConfig(filePath string) *Configure {
	var cfg Configure

	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		return DefaultConfig()
	}
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return DefaultConfig()
	}
	log.Printf("config: %+v", cfg)
	return &cfg
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
	if config.Monitor.Path != "" {
		return string(config.Monitor.Path)
	}
	return HomePath(".ipfs")
}

// IpfsClusterPath ...
func IpfsClusterPath() string {
	if config.Monitor.ClusterPath != "" {
		return string(config.Monitor.ClusterPath)
	}
	return HomePath(".ipfs-cluster")
}

// DefaultConfig ...
func DefaultConfig() *Configure {
	return &Configure{}
}

// Config ...
func Config() *Configure {
	if config == nil {
		panic("nil config")
	}
	return config
}

// DefaultString ...
func DefaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

// CheckExist ...
func (cfg *MonitorProperty) CheckExist() bool {
	_, err := os.Stat(filepath.Join(cfg.RootPath, DefaultFileName))
	if err != nil {
		errors.ErrorStack(err)
		return false
	}
	return true
}

// DefaultMonitorProperty ...
func DefaultMonitorProperty() *MonitorProperty {
	return &MonitorProperty{
		Version:             "v0",
		CommandName:         "ipfs",
		ServiceCommandName:  "ipfs-cluster-service",
		RootPath:            "",
		HostType:            "",
		RemoteIP:            "127.0.0.1",
		RemotePort:          ":7758",
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
