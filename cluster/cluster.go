package cluster

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"log"
	"os"
	"os/exec"
	"time"
)

// MonitorEnviron ...
var monitorEnviron []string

func init() {
	monitorEnviron = api.Config().MonitorEnviron
}

func firstRunIPFS() {
	cmd := exec.Command("ipfs", "init")
	cmd.Env = os.Environ()
	log.Println(cmd.Env)
	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes), err)
}

func firstRunService() {

}

// WaitingForInitialize ...
func WaitingForInitialize(ctx context.Context) bool {
	for {
		if !api.IsInitialized() {
			time.Sleep(time.Second * 5)
			select {
			case <-ctx.Done():
				return false
			default:
				continue
			}
		}
		return true
	}
}

// Run ...
func Run(ctx context.Context) {
	if WaitingForInitialize(ctx) {

	}
}

// InitRun ...
func InitRun() {

}
