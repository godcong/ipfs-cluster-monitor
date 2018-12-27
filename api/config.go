package api

import (
	"fmt"
	"github.com/juju/errors"
	"os"
	"os/user"
	"path/filepath"
)

// DefaultFileName ...
const DefaultFileName = "monitor.json"

var isInitialized bool

// ClassType ...
type ClassType string

// ClassHost ...
var (
	ClassHost   ClassType = "host"
	ClassClient ClassType = "client"
)

// Config ...
type Config struct {
	RootPath string
	Secret   string
	Class    string
	RemoteIP string
}

var config *Config

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
	}
}

// InitLoader ...
func (cfg *Config) InitLoader() {
	fileInfo, err := os.Stat(cfg.RootPath + "/" + DefaultFileName)
	if err != nil {
		err := os.MkdirAll(cfg.RootPath, os.ModePerm)
		if err != nil {
			errors.ErrorStack(err)
			panic(err)
		}
	}
	isInitialized = true
}

// Make ...
func (cfg *Config) Make() *Config {

}

// IsInitialized ...
func IsInitialized() bool {
	return isInitialized
}
