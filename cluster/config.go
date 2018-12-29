package cluster

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
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

// InitService ...
const InitService = "service"

var isInitialized bool

// HostType ...
type HostType string

// ClassHost ...
var (
	HostServer HostType = "server"
	HostClient HostType = "client"
)

type EnvironIPFS string
type EnvironSecret string
type EnvironService string

func (s EnvironIPFS) String() string {
	return strings.Join([]string{"IPFS_PATH", string(s)}, "=")
}

func (s EnvironSecret) String() string {
	return strings.Join([]string{"IPFS_CLUSTER_SECRET", string(s)}, "=")
}

func (s EnvironService) String() string {
	return strings.Join([]string{"IPFS_CLUSTER_PATH", string(s)}, "=")
}

type Environ struct {
	Ipfs    EnvironIPFS
	Secret  EnvironSecret
	Service EnvironService
}

// Configuration ...
type Configuration struct {
	Version             string
	RootPath            string
	CommandName         string
	ServiceCommandName  string
	Secret              string
	HostType            HostType
	RemoteIP            string
	RemotePort          string
	Environ             Environ
	Interval            time.Duration
	ServerCheckInterval time.Duration
	MonitorInterval     time.Duration
	ResetWaiting        int
}

var cfg *Configuration

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
	cfg = DefaultConfig()
	cfg.InitLoader()
}

func defaultPath(name string) string {
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

// Config ...
func Config() *Configuration {
	return cfg
}

func getClusterPath() string {
	rootPath, b := os.LookupEnv("IPFS_CLUSTER_MONITOR")
	if !b {
		rootPath = defaultPath(".ipfs-cluster-monitor")
	}
	return rootPath
}

// DefaultConfig ...
func DefaultConfig() *Configuration {
	return &Configuration{
		Version:             "v0",
		CommandName:         "ipfs",
		ServiceCommandName:  "ipfs-cluster-service",
		RootPath:            getClusterPath(),
		HostType:            HostServer,
		RemoteIP:            "127.0.0.1",
		RemotePort:          ":7758",
		Interval:            1 * time.Second,
		MonitorInterval:     5 * time.Second,
		ServerCheckInterval: 60 * time.Second,
		ResetWaiting:        30,
	}
}

// SetClient ...
func (cfg *Configuration) SetClient(remoteIP string) {
	cfg.HostType = HostClient
	cfg.RemoteIP = remoteIP
}

// InitLoader ...
func (cfg *Configuration) InitLoader() {

	if !IsInitialized() {
		return
	}

	file, err := os.OpenFile(filepath.Join(cfg.RootPath, DefaultFileName), os.O_RDONLY|os.O_SYNC, os.ModePerm)
	CheckError(err)
	defer file.Close()

	dec := jsoniter.NewDecoder(file)

	err = dec.Decode(cfg)

	CheckError(err)
}

// CheckExist ...
func (cfg *Configuration) CheckExist() bool {
	_, err := os.Stat(filepath.Join(cfg.RootPath, DefaultFileName))
	if err != nil {
		errors.ErrorStack(err)
		return false
	}
	return true
}

// Marshal ...
func (cfg *Configuration) Marshal() ([]byte, error) {
	return jsoniter.Marshal(cfg)
}

// SetEnv ...
func (cfg *Configuration) SetEnv(value interface{}) {
	if value != nil {
		switch v := value.(type) {
		case EnvironIPFS:
			cfg.Environ.Ipfs = v
		case EnvironSecret:
			cfg.Environ.Secret = v
		case EnvironService:
			cfg.Environ.Service = v
		}
	}
}

// GetEnv ...
func (cfg *Configuration) GetEnv() []string {
	env := os.Environ()
	if cfg.Environ.Service != "" {
		env = append(env, cfg.Environ.Service.String())
	}
	if cfg.Environ.Secret != "" {
		env = append(env, cfg.Environ.Secret.String())
	}
	if cfg.Environ.Ipfs != "" {
		env = append(env, cfg.Environ.Ipfs.String())
	}
	return env
}

// Make ...
func (cfg *Configuration) Make() {
	err := os.Chdir(cfg.RootPath)
	if err != nil {
		err := os.MkdirAll(cfg.RootPath, os.ModePerm)
		CheckError(err)
	}

	file, err := os.OpenFile(filepath.Join(cfg.RootPath, DefaultFileName), os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	log.Println("created:", file.Name())
	CheckError(err)
	defer file.Close()

	enc := jsoniter.NewEncoder(file)
	err = enc.Encode(*cfg)
	CheckError(err)

	cfile, err := os.Create(filepath.Join(cfg.RootPath, InitIPFS))
	log.Println("created:", cfile.Name())
	CheckError(err)
	defer cfile.Close()

	sfile, err := os.Create(filepath.Join(cfg.RootPath, InitService))
	log.Println("created:", sfile.Name())
	CheckError(err)
	defer sfile.Close()

}

// IsInitialized ...
func IsInitialized() bool {
	if isInitialized == false {
		isInitialized = cfg.CheckExist()
	}

	return isInitialized
}

// CheckError ...
func CheckError(err error) {
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
}
