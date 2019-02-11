package config

import (
	"fmt"
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
	Enable      bool   `toml:"enable"`
	Type        string `toml:"type"`
	AddrPort    string `toml:"addr_port"`
	Workspace   string `toml:"workspace"`
	Secret      string `toml:"secret"`
	Bootstrap   string `toml:"bootstrap"`
	IpfsPath    string `toml:"ipfs_path"`
	ClusterPath string `toml:"cluster_path"`
}

// MustMonitor ...
func MustMonitor(secret, boot, workspace string) *Monitor {
	return &Monitor{
		Secret:      DefaultString(secret, "27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34"),
		Bootstrap:   DefaultString(boot, "/ip4/47.101.169.94/tcp/9096/ipfs/QmU58AYMghsHEMq6gSrLNT1kVPigG3gpvfaifeUuXKXeLs"),
		IpfsPath:    DefaultString(filepath.Join(workspace, Ipfs), HomePath(".ipfs")),
		ClusterPath: DefaultString(filepath.Join(workspace, Cluster), HomePath(".ipfs-cluster")),
	}
}

// Env ...
func (m *Monitor) Env() (env []string) {
	env = os.Environ()
	env = append(env, strings.Join([]string{"IPFS_PATH", string(m.IpfsPath)}, "="))
	env = append(env, strings.Join([]string{"CLUSTER_SECRET", string(m.Secret)}, "="))
	env = append(env, strings.Join([]string{"IPFS_CLUSTER_PATH", string(m.ClusterPath)}, "="))

	log.Println(env)
	return
}

// MonitorProperty ...
type MonitorProperty struct {
	Version             string        `toml:"version"`
	IpfsCommandName     string        `toml:"ipfs_command_name"`
	ClusterCommandName  string        `toml:"cluster_command_name"`
	RemoteIP            string        `toml:"remote_ip"`
	RemotePort          string        `toml:"remote_port"`
	Interval            time.Duration `toml:"interval"`
	ServerCheckInterval time.Duration `toml:"server_check_interval"`
	MonitorInterval     time.Duration `toml:"monitor_interval"`
	ResetWaiting        int           `toml:"reset_waiting"`
}

// Configure ...
type Configure struct {
	Initialize      bool            `toml:"-"`
	Root            string          `toml:"-"`
	ConfigName      string          `toml:"-"`
	Monitor         Monitor         `toml:"monitor"`
	MonitorProperty MonitorProperty `toml:"monitor_property"`
	//Database        Database        `toml:"database"`
	//Censor          HostInfo        `toml:"censor"`
	//Node            HostInfo        `toml:"node"`
	//Media           Media           `toml:"media"`
	//Queue           Queue           `toml:"queue"`
	GRPC GRPC `toml:"grpc"`
	//REST            REST            `toml:"rest"`
	IPFS IPFS `toml:"ipfs"`

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
	dir, name := filepath.Split(filePath[0])
	config.Root = dir
	config.ConfigName = name
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
	cfg.Initialize = true
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
			Enable:      false,
			Type:        "tcp",
			AddrPort:    "localhost:7784",
			Secret:      "27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34",
			Bootstrap:   "/ip4/47.101.169.94/tcp/9096/ipfs/QmU58AYMghsHEMq6gSrLNT1kVPigG3gpvfaifeUuXKXeLs",
			IpfsPath:    HomePath(".ipfs"),
			ClusterPath: HomePath(".ipfs-cluster"),
		},
		MonitorProperty: MonitorProperty{
			Version:             "",
			IpfsCommandName:     "ipfs",
			ClusterCommandName:  "ipfs-cluster-service",
			RemoteIP:            "",
			RemotePort:          "",
			Interval:            3 * time.Second,
			ServerCheckInterval: 3 * time.Second,
			MonitorInterval:     3 * time.Second,
			ResetWaiting:        0,
		},
		GRPC:      GRPC{},
		IPFS:      IPFS{},
		Requester: Requester{},
		Callback:  Callback{},
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
