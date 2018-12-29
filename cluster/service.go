package cluster

import (
	"context"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
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

// runService ...
func runService(ctx context.Context) {
	if useBootstrap() {
		boot := getServiceBootstrap()
		if boot != "" {
			go optimizeRunCMD(cfg.ServiceCommandName, "daemon", "--bootstrap", boot)
			return
		}
	}
	go optimizeRunCMD(cfg.ServiceCommandName, "daemon")
}

// DeletePeers ...
func DeletePeers(peerID string) error {
	request, err := http.NewRequest(http.MethodDelete, "http://localhost:9094/peers/"+peerID, nil)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	log.Println(string(bytes))
	return nil
}
