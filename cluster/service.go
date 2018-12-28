package cluster

import (
	"github.com/juju/errors"
	"log"
	"os"
	"os/exec"
)

func firstRunService() {
	cmd := exec.Command("ipfs-cluster-service", "init")
	cmd.Env = os.Environ()
	if clusterEnviron != nil {
		cmd.Env = append(cmd.Env, clusterEnviron...)
	}

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
}
