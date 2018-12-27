package api

import (
	"fmt"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

// DefaultFileName ...
const DefaultFileName = "monitor.json"

var isInitialized bool

// HostType ...
type HostType string

// ClassHost ...
var (
	HostServer HostType = "server"
	HostClient HostType = "client"
)

// Config ...
type Config struct {
	RootPath string
	Secret   string
	HostType HostType
	RemoteIP string
}

var config *Config

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate)
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

// DefaultConfig ...
func DefaultConfig() *Config {
	rootPath, b := os.LookupEnv("IPFS_CLUSTER_MONITOR")
	if !b {
		rootPath = defaultPath(".ipfs_cluster_monitor")
	}

	return &Config{
		RootPath: rootPath,
		HostType: HostServer,
		RemoteIP: "127.0.0.1",
	}
}

// SetClient ...
func (cfg *Config) SetClient(remoteIP string) {
	cfg.HostType = HostClient
	cfg.RemoteIP = remoteIP
}

// InitLoader ...
func (cfg *Config) InitLoader() {
	if !IsInitialized(cfg) {
		cfg.Make()
		isInitialized = true
	}

	file, err := os.OpenFile(cfg.RootPath+"/"+DefaultFileName, os.O_RDONLY|os.O_SYNC, os.ModePerm)

	CheckError(err)

	dec := jsoniter.NewDecoder(file)

	err = dec.Decode(cfg)
	CheckError(err)
}

// CheckExist ...
func (cfg *Config) CheckExist() bool {
	_, err := os.Stat(cfg.RootPath + "/" + DefaultFileName)
	if err != nil {
		errors.ErrorStack(err)
		return false
	}
	return true
}

// Make ...
func (cfg *Config) Make() {
	err := os.Chdir(cfg.RootPath)
	if err != nil {
		err := os.MkdirAll(cfg.RootPath, os.ModePerm)
		CheckError(err)
	}

	file, err := os.OpenFile(cfg.RootPath+"/"+DefaultFileName, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	CheckError(err)

	enc := jsoniter.NewEncoder(file)

	err = enc.Encode(*cfg)
	CheckError(err)
}

// IsInitialized ...
func IsInitialized(cfg *Config) bool {
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