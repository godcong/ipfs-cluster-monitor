package cluster

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"github.com/juju/errors"
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
	if monitorEnviron != nil {
		cmd.Env = append(cmd.Env, monitorEnviron...)
	}

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}

	//log.Println(err)
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
