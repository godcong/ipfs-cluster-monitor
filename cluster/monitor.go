package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// MonitorPeers ...
const MonitorPeers = "peers"

// Peer ...
type Peer struct {
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

var monitor sync.Map

// runMonitor ...
func runMonitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			peers, err := getPeers()
			if err != nil {
				monitor.Delete(MonitorPeers)
				time.Sleep(cfg.MonitorInterval)
				errors.ErrorStack(err)
				log.Println(err)
				continue
			}
			monitor.Store(MonitorPeers, peers)
			time.Sleep(cfg.MonitorInterval)
		}
		//get info
	}
}

func getPeers() ([]Peer, error) {
	response, err := http.Get("http://localhost:9094/peers")
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var peers []Peer

	err = jsoniter.Unmarshal(bytes, &peers)
	if err != nil {
		return nil, err
	}

	//monitor.Store("peers", peers)

	return peers, nil
}

// GetPeers ...
func GetPeers() []Peer {
	if peers, b := monitor.Load(MonitorPeers); b {
		return peers.([]Peer)
	}
	return nil
}
