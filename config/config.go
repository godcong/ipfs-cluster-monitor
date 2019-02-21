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

// ClusterClient ...
type ClusterClient struct {
	Secret      string   `toml:"secret"`
	Bootstrap   []string `toml:"bootstrap"`
	ClusterPath string   `toml:"cluster_path"`
}

// IPFSClient ...
type IPFSClient struct {
	IpfsPath string `toml:"ipfs_path"`
}

// Monitor ...
type Monitor struct {
	Token         string         `toml:"token"`
	Enable        bool           `toml:"enable"`
	Type          string         `toml:"type"`
	Addr          string         `toml:"addr"`
	Port          string         `toml:"port"`
	Workspace     string         `toml:"workspace"`
	IPFSClient    *IPFSClient    `toml:"ipfs_client"`
	ClusterClient *ClusterClient `toml:"cluster_client"`
}

// MustIPFSClient ...
func MustIPFSClient(ws string) *IPFSClient {
	return &IPFSClient{
		IpfsPath: DefaultString(filepath.Join(ws, "data", Ipfs), HomePath(".ipfs")),
	}
}

// MustClusterClient ...
func MustClusterClient(ws string, sec string, boot string) *ClusterClient {
	return &ClusterClient{
		Secret:      sec,
		Bootstrap:   []string{boot},
		ClusterPath: DefaultString(filepath.Join(ws, "data", Cluster), HomePath(".ipfs-cluster")),
	}
}

// MustMonitor ...
func MustMonitor(secret, boot, workspace string) *Monitor {
	var ipfs IPFSClient
	var cluster ClusterClient
	if workspace != "" {
		&ipfs = MustIPFSClient(workspace)
		&cluster = MustClusterClient(workspace, secret, boot)
	}
	log.Debug(workspace, ipfs, cluster)
	return &Monitor{
		Enable:        true,
		Workspace:     workspace,
		IPFSClient:    &ipfs,
		ClusterClient: &cluster,
	}
}

// Env ...
func (m *Monitor) Env() (env []string) {
	var e error
	e = os.Setenv("IPFS_PATH", string(m.IPFSClient.IpfsPath))
	if e != nil {
		panic(e)
	}

	e = os.Setenv("CLUSTER_SECRET", string(m.ClusterClient.Secret))
	if e != nil {
		panic(e)
	}

	e = os.Setenv("IPFS_CLUSTER_PATH", string(m.ClusterClient.ClusterPath))
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
	Interval            time.Duration `toml:"interval"`
	ServerCheckInterval time.Duration `toml:"server_check_interval"`
	MonitorInterval     time.Duration `toml:"monitor_interval"`
	ResetWaiting        int           `toml:"reset_waiting"`
}

// Configure ...
type Configure struct {
	Mode            int             `toml:"mode"`
	Initialize      bool            `toml:"-"`
	ConfigPath      string          `toml:"-"`
	RunPath         string          `toml:"-"` //运行路径(启动加载)
	ConfigName      string          `toml:"-"` //配置文件名
	ClusterServer   ClusterClient   `toml:"cluster_server"`
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
func Initialize(runPath string, configPath ...string) error {
	log.Debug(runPath, configPath)
	config = DefaultConfig(runPath)
	config.LoadConfig(configPath[0])
	config.ConfigPath, config.ConfigName = filepath.Split(configPath[0])
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
func (obj *Configure) LoadConfig(filePath string) *Configure {
	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	if err != nil {
		log.Error("config open:", err)
		return obj
	}
	defer openFile.Close()
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(obj)
	if err != nil {
		log.Error("config decode:", err)
		return obj
	}
	obj.Initialize = true
	log.Debugf("config: %+v", obj)
	return obj
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
	if config.Monitor.IPFSClient.IpfsPath != "" {
		return config.Monitor.IPFSClient.IpfsPath
	}
	return HomePath(".ipfs")
}

// ClusterPath ...
func ClusterPath() string {
	if config.Monitor.ClusterClient.ClusterPath != "" {
		return config.Monitor.ClusterClient.ClusterPath
	}
	return HomePath(".ipfs-cluster")
}

// DefaultConfig ...
func DefaultConfig(runPath string) *Configure {
	return &Configure{
		Initialize: false,
		RunPath:    runPath,
		Monitor: Monitor{
			Token:      "2UQEoTCGV7j689CEImQDAjcv7k1X0ZpxT2yzCX8vaqRg1vKp5f0uScvPVB7yuZPP",
			Enable:     true,
			Type:       "tcp",
			Addr:       "localhost",
			Port:       ":7784",
			Workspace:  runPath,
			IPFSClient: MustIPFSClient(runPath),
			ClusterClient: MustClusterClient(runPath,
				"27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34",
				"/ip4/47.101.169.94/tcp/9096/ipfs/QmeQzPKd7HzKZwBKNmnJnyub3YyCBvtcWraaJKEKk1BWmx"),
		},
		MonitorProperty: *DefaultMonitorProperty(runPath),
		GRPC: GRPC{
			Enable: true,
			Type:   "",
			Path:   "",
			Port:   "",
		},
		REST: REST{},
		IPFS: IPFS{
			Host: "",
			Port: "",
		},
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

// SetRoot ...
func SetRunPath(fp string) (err error) {
	config.RunPath, err = filepath.Abs(filepath.Dir(fp)) //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
		config.RunPath = ""
	}
	//TODO:Maybe move to together
	if config.ConfigPath == "" {
		config.ConfigPath = config.RunPath
	}
	return
}

// FD config file dictionary
func (c *Configure) FD() string {
	return filepath.Join(c.ConfigPath, c.ConfigName)
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
	_, err := os.Stat(filepath.Join(c.ConfigPath, DefaultFileName))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// DefaultMonitorProperty ...
func DefaultMonitorProperty(runPath string) *MonitorProperty {
	return &MonitorProperty{
		Version:             "v0",
		IpfsCommandName:     filepath.Join(runPath, "ipfs"),
		ClusterCommandName:  filepath.Join(runPath, "ipfs-cluster-service"),
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
