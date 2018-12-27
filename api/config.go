package api

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	RootPath string
	FileName string
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

func DefaultConfig() *Config {
	rootPath, b := os.LookupEnv("IPFS_CLUSTER_MONITOR")
	if !b {
		rootPath = defaultPath(".ipfs_cluster_monitor")
	}

	chdir := os.Chdir(rootPath)
	log.Println(chdir)

	return &Config{
		RootPath: rootPath,
		FileName: "monitor.json",
	}
}
