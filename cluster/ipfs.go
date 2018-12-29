package cluster

import (
	"context"
	"github.com/juju/errors"
	"log"
	"os/exec"
)

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
