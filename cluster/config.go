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

// Configuration ...
type Configuration struct {
	Version            string
	RootPath           string
	CommandName        string
	ServiceCommandName string
	Secret             string
	HostType           HostType
	RemoteIP           string
	ClusterEnviron     []string
	MonitorInterval    time.Duration
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

// DefaultConfig ...
func DefaultConfig() *Configuration {
	rootPath, b := os.LookupEnv("IPFS_CLUSTER_MONITOR")
	if !b {
		rootPath = defaultPath(".ipfs-cluster-monitor")
	}

	return &Configuration{
		Version:            "v0",
		CommandName:        "ipfs",
		ServiceCommandName: "ipfs-cluster-service",
		RootPath:           rootPath,
		HostType:           HostServer,
		RemoteIP:           "127.0.0.1",
		MonitorInterval:    5 * time.Second,
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

	file, err := os.OpenFile(cfg.RootPath+"/"+DefaultFileName, os.O_RDONLY|os.O_SYNC, os.ModePerm)
	CheckError(err)
	defer file.Close()

	dec := jsoniter.NewDecoder(file)

	err = dec.Decode(cfg)

	CheckError(err)
}

// CheckExist ...
func (cfg *Configuration) CheckExist() bool {
	_, err := os.Stat(cfg.RootPath + "/" + DefaultFileName)
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
func (cfg *Configuration) SetEnv(key, value string) {
	if value != "" {
		cfg.ClusterEnviron = append(cfg.ClusterEnviron, strings.Join([]string{key, value}, "="))
	}
}

// GetEnv ...
func (cfg *Configuration) GetEnv() []string {
	if cfg.ClusterEnviron != nil {
		return append(os.Environ(), cfg.ClusterEnviron...)
	}
	return os.Environ()
}

// Make ...
func (cfg *Configuration) Make() {
	err := os.Chdir(cfg.RootPath)
	if err != nil {
		err := os.MkdirAll(cfg.RootPath, os.ModePerm)
		CheckError(err)
	}

	file, err := os.OpenFile(cfg.RootPath+"/"+DefaultFileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	CheckError(err)
	defer file.Close()

	enc := jsoniter.NewEncoder(file)
	err = enc.Encode(*cfg)
	CheckError(err)

	cfile, err := os.Create(cfg.RootPath + "/" + InitIPFS)
	CheckError(err)
	defer cfile.Close()

	sfile, err := os.Create(cfg.RootPath + "/" + InitService)
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
