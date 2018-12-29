package cluster

import (
	"context"
	"github.com/juju/errors"
	"log"
	"os/exec"
)

func firstRunService() {
	cmd := exec.Command(cfg.ServiceCommandName, "init")
	cmd.Env = cfg.GetEnv()

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
}

func optimizationFirstRunService(ctx context.Context) {
	err := optimizeRunCMD(cfg.ServiceCommandName, "init")
	if err != nil {
		panic(err)
	}
}
