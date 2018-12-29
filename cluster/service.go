package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
)

// ServicePeer ...
type ServicePeer struct {
	ID                    string        `json:"id"`
	Addresses             []string      `json:"addresses"`
	ClusterPeers          []string      `json:"cluster_peers"`
	ClusterPeersAddresses []interface{} `json:"cluster_peers_addresses"`
	Version               string        `json:"version"`
	Commit                string        `json:"commit"`
	RPCProtocolVersion    string        `json:"rpc_protocol_version"`
	Error                 string        `json:"error"`
	Ipfs                  struct {
		ID        string   `json:"id"`
		Addresses []string `json:"addresses"`
		Error     string   `json:"error"`
	} `json:"ipfs"`
	Peername string `json:"peername"`
}

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
	if isClient() {
		boot := getServiceBootstrap()
		if boot != "" {
			go optimizeRunCMD(cfg.ServiceCommandName, "daemon", "--bootstrap", boot)
			return
		}
	}
	go optimizeRunCMD(cfg.ServiceCommandName, "daemon")
}

func getPeers() ([]ServicePeer, error) {
	response, err := http.Get("http://localhost:9094/peers")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var peers []ServicePeer

	err = jsoniter.Unmarshal(bytes, &peers)
	if err != nil {
		return nil, err
	}

	//monitor.Store("peers", peers)

	return peers, nil
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

// GetPeers ...
func GetPeers() []ServicePeer {
	if peers, b := monitor.Load(MonitorPeers); b {
		return peers.([]ServicePeer)
	}
	return nil
}
