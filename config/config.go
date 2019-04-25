package config

import (
	"github.com/godcong/ipfs-cluster-monitor/proto"
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

// IpfsTmp ...
const IpfsTmp = "ipfs_tmp"

// Cluster ...
const Cluster = "cluster"

// ClusterTmp ...
const ClusterTmp = "cluster_tmp"

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
	Type string `toml:"type"`
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

// ClusterClient ...
type ClusterClient struct {
	Secret      string   `toml:"secret"`
	Bootstrap   []string `toml:"bootstrap"`
	ClusterPath string   `toml:"cluster_path"`
}

// IPFSClient ...
type IPFSClient struct {
	IpfsPath  string `toml:"ipfs_path"`
	Bootstrap string `toml:"bootstrap"`
}

// MonitorServer ...
type MonitorServer struct {
}

// Monitor ...
type Monitor struct {
	GRPC          GRPC            `toml:"grpc"`
	REST          REST            `toml:"rest"`
	Mode          proto.StartMode `toml:"mode"`
	Token         string          `toml:"token"`
	Enable        bool            `toml:"enable"`
	Workspace     string          `toml:"workspace"`
	IPFSClient    IPFSClient      `toml:"ipfs_client"`
	ClusterClient ClusterClient   `toml:"cluster_client"`
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
	Workspace           string        `toml:"workspace"`
}

// Logger ...
type Logger struct {
	Path  string `toml:"path"`
	Level string `toml:"level"`
}

// Custom ...
type Custom struct {
	MaxSize   string `toml:"max_size"`
	Workspace string `toml:"workspace"`
}

// Configure ...
type Configure struct {
	UseCustom       bool            `toml:"use_custom"`
	Mode            int             `toml:"mode"`
	Logger          Logger          `toml:"logger"`
	Initialize      bool            `toml:"-"`
	ConfigPath      string          `toml:"-"`
	RunPath         string          `toml:"-"` //运行路径(启动加载)
	ConfigName      string          `toml:"-"` //配置文件名
	Monitor         Monitor         `toml:"monitor"`
	MonitorProperty MonitorProperty `toml:"monitor_property"`
	GRPC            GRPC            `toml:"grpc"`
	REST            REST            `toml:"rest"`
	IPFS            IPFS            `toml:"ipfs"`
	Custom          Custom          `toml:"custom"`
	Callback        Callback        `toml:"callback"`
}

var config *Configure

// MustIPFSClient ...
func MustIPFSClient(ws string) *IPFSClient {
	return &IPFSClient{
		IpfsPath:  DefaultString(filepath.Join(ws, "data", Ipfs), HomePath(".ipfs")),
		Bootstrap: DefaultString("", "/ip4/47.101.169.94/tcp/4001/ipfs/QmRkiD3iWg3W2yJG5g9YD3jhRn5CX9HqG2uVAGjNFdNGQ2"),
	}
}

// MustClusterClient ...
func MustClusterClient(ws string, sec string, boot string) *ClusterClient {
	sec = DefaultString(sec, "27b3f5c4e330c069cc045307152345cc391cb40e6dcabf01f98ae9cdc9dabb34")
	boot = DefaultString(boot, "/ip4/47.101.169.94/tcp/9096/ipfs/QmeQzPKd7HzKZwBKNmnJnyub3YyCBvtcWraaJKEKk1BWmx")
	return &ClusterClient{
		Secret:      sec,
		Bootstrap:   []string{boot},
		ClusterPath: DefaultString(filepath.Join(ws, "data", Cluster), HomePath(".ipfs-cluster")),
	}
}

// MustMonitor ...
func MustMonitor(mode proto.StartMode, secret, boot, workspace string) *Monitor {
	var ipfs IPFSClient
	var cluster ClusterClient
	if workspace != "" {
		ipfs = *MustIPFSClient(workspace)
		cluster = *MustClusterClient(workspace, secret, boot)
	}
	log.Debug(workspace, ipfs, cluster)
	return &Monitor{
		GRPC: GRPC{
			Enable: true,
			Type:   "tcp",
			Path:   "47.101.169.94",
			Port:   ":7774",
		},
		REST:          REST{},
		Mode:          mode,
		Token:         "",
		Enable:        true,
		Workspace:     workspace,
		IPFSClient:    ipfs,
		ClusterClient: cluster,
	}
}

// Env ...
func (m *Monitor) Env() (env []string) {
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()
	p := os.Getenv("PATH")
	if err := os.Setenv("PATH", p+":"+m.Workspace); err != nil {
		err = xerrors.Errorf("PATH error:%+v", err)
		return nil
	}

	if err := os.Setenv("IPFS_PATH", m.IPFSClient.IpfsPath); err != nil {
		err = xerrors.Errorf("IPFS_PATH error:%+v", err)
		return nil
	}

	if m.Mode == proto.StartMode_Cluster {
		if err := os.Setenv("CLUSTER_SECRET", m.ClusterClient.Secret); err != nil {
			err = xerrors.Errorf("CLUSTER_SECRET error:%+v", err)
			return nil
		}
		if err := os.Setenv("IPFS_CLUSTER_PATH", m.ClusterClient.ClusterPath); err != nil {
			err = xerrors.Errorf("IPFS_CLUSTER_PATH error:%+v", err)
			return nil
		}
	}

	env = os.Environ()
	return
}

// Initialize ...
func Initialize(runPath string, configPath ...string) error {
	log.Info(runPath, configPath)

	dir, file := filepath.Split(configPath[0])
	if file == "" {
		file = "config.toml"
	}

	s, e := filepath.Abs(dir)
	if e != nil {
		s = ""
	}
	log.Info(s)
	config = DefaultConfig(dir)

	e = config.ReadConfig(configPath[0])
	config.RunPath = runPath
	config.ConfigPath, config.ConfigName = dir, file
	if e != nil || !config.Initialize {
		return create(config)
	}
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

// ReadConfig ...
func (c *Configure) ReadConfig(filePath string) (e error) {
	openFile, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, os.ModePerm)

	if err != nil {
		return xerrors.Errorf("open config file error:%+v", err)
	}

	defer openFile.Close()
	decoder := toml.NewDecoder(openFile)
	err = decoder.Decode(c)
	if err != nil {
		return xerrors.Errorf("decode config file error:%+v", err)
	}
	c.Initialize = true
	return nil
}

func create(cfg *Configure) error {
	//file, e := os.OpenFile(configure.ConfigPath, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	_, e := os.Stat(cfg.ConfigPath)
	log.Info("init dictionary stat:", e)
	if os.IsNotExist(e) {
		log.Println("dictionary not exist creating... ")
		_ = os.MkdirAll(cfg.ConfigPath, os.ModePerm)
	}
	file, e := os.OpenFile(config.FileConfig(), os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if e != nil {
		return xerrors.Errorf("make file error:%+v", e)
	}
	defer file.Close()
	//create config file
	enc := toml.NewEncoder(file)
	e = enc.Encode(*cfg)
	if e != nil {
		return xerrors.Errorf("encode file:%w", e)
	}
	log.Println("created:", file.Name())
	//create directory
	e = os.MkdirAll(filepath.Join(cfg.Monitor.Workspace, "data"), os.ModePerm)
	if e != nil {
		log.Println("make workspace err:", cfg.Monitor.Workspace, e)
	}
	//
	cfile, e := os.Create(filepath.Join(cfg.Monitor.Workspace, IpfsTmp))
	if e != nil {
		log.Println(e)
		return xerrors.Errorf("ipfs file:%w", e)
	}
	log.Println("created:", cfile.Name())
	defer cfile.Close()

	sfile, e := os.Create(filepath.Join(cfg.Monitor.Workspace, ClusterTmp))
	if e != nil {
		log.Println(e)
		return xerrors.Errorf("cluster file:%w", e)
	}
	log.Println("created:", sfile.Name())
	defer sfile.Close()

	time.Sleep(3 * time.Second)
	cfg.Initialize = true

	return nil
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
			return filepath.Join(home, name)
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
func DefaultConfig(ws string) *Configure {
	return &Configure{
		Initialize:      false,
		Monitor:         *MustMonitor(proto.StartMode_Cluster, "", "", ws),
		MonitorProperty: *MustMonitorProperty(ws),
		GRPC: GRPC{
			Enable: true,
			Type:   "tcp",
			Path:   "127.0.0.1",
			Port:   ":7784",
		},
		REST:     REST{},
		IPFS:     IPFS{},
		Callback: Callback{},
	}
}

// Config ...
func Config() *Configure {
	if config == nil {
		panic("nil config")
	}
	return config
}

// SetRunPath ...
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

// FileConfig config file dictionary
func (c *Configure) FileConfig() string {
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

// MustMonitorProperty ...
func MustMonitorProperty(ws string) *MonitorProperty {
	return &MonitorProperty{
		Version:             "v0",
		IpfsCommandName:     "ipfs",
		ClusterCommandName:  "ipfs-cluster-service",
		Interval:            1 * time.Second,
		ServerCheckInterval: 60 * time.Second,
		MonitorInterval:     5 * time.Second,
		ResetWaiting:        30,
		Workspace:           ws,
	}
}

// SetMonitor ...
func SetMonitor(monitor *Monitor) {
	if config != nil {
		config.Monitor = *monitor
	}
}

// CheckError ...
func CheckError(err error) error {
	return xerrors.Errorf("check err:%w", err)
}
