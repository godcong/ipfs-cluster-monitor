package cluster

import (
	"context"
	"fmt"
	"github.com/juju/errors"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

// IPFSInfo ...
type IPFSInfo struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

// runIPFS ...
func runIPFS(ctx context.Context) {
	go optimizeRunCMD("ipfs", "daemon")
}

func firstRunIPFS() {
	cmd := exec.Command(cfg.CommandName, "init")
	cmd.Env = cfg.GetEnv()

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
}

func optimizationFirstRunIPFS(ctx context.Context) {
	err := optimizeRunCMD(cfg.CommandName, "init")
	if err != nil {
		panic(err)
	}
}

func defaultIPFS() string {
	home := os.Getenv("HOME")
	if home == "" {
		usr, err := user.Current()
		if err != nil {
			panic(fmt.Sprintf("cannot get current user: %s", err))
		}
		home = usr.HomeDir
	}

	return filepath.Join(home, ".ipfs")
}
