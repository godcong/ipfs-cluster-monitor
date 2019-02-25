package cluster

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"time"
)

// IpfsInfo ...
type IpfsInfo struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

// IpfsAddress ...
type IpfsAddress struct {
	Addresses []string `json:"addresses"`
}

// IpfsList ...
type IpfsList struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Detail  IpfsAddress
}

// RunIPFS ...
func RunIPFS(ctx context.Context, cfg *config.Configure) {
	go optimizeRunCMD(ctx, cfg.MonitorProperty.IpfsCommandName, cfg.Monitor.Env(), "daemon")
}

// RunIPFSInit ...
func RunIPFSInit(ctx context.Context, cfg *config.Configure) error {
	cmd := exec.CommandContext(ctx, cfg.MonitorProperty.IpfsCommandName, "init")
	cmd.Env = cfg.Monitor.Env()

	bytes, err := cmd.CombinedOutput()
	log.Println(string(bytes))
	if err != nil {
		return xerrors.Errorf("first run ipfs:%w", err)
	}
	return nil
}

func optimizationFirstRunIPFS(ctx context.Context, cfg *config.Configure) {
	err := optimizeRunCMD(ctx, cfg.MonitorProperty.IpfsCommandName, cfg.Monitor.Env(), "init")
	if err != nil {
		panic(err)
	}
}

// GetIpfsInfo ...
func GetIpfsInfo() (*IpfsInfo, error) {
	return getIpfsInfo()
}

func getIpfsInfo() (*IpfsInfo, error) {
	var ipfs IpfsInfo
	response, err := http.Get("http://localhost:5001/api/v0/id")
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(bytes, &ipfs)
	if err != nil {
		return nil, err
	}
	return &ipfs, nil
}

func getRemoteIpfsList() (*IpfsList, error) {
	var ipfs IpfsList
	//TODO:remote ip address
	response, err := http.Get("")
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = jsoniter.Unmarshal(bytes, &ipfs)
	if err != nil {
		return nil, err
	}
	return &ipfs, nil
}

// WaitingIPFS ...
func WaitingIPFS(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			log.Println("ipfs done")
			return
		default:
			log.Println("waiting ipfs")
			time.Sleep(1 * time.Second)
			_, err = getIpfsInfo()
			if err == nil {
				return
			}
		}
	}
}

// Pin ...
func Pin(arg string) error {
	v := url.Values{
		"arg": {arg},
	}
	resp, e := http.Get("http://localhost:5001/api/v0/pin/add?" + v.Encode())
	if e != nil {
		log.Error(e)
	}
	bytes, e := ioutil.ReadAll(resp.Body)
	log.Info("pin res:", string(bytes), e)
	return e
}

// SwarmAddress ...
func SwarmAddress(arg string) error {
	v := url.Values{
		"arg": {arg},
	}
	resp, e := http.Get("http://localhost:5001/api/v0/swarm/connect?" + v.Encode())
	if e != nil {
		log.Error(e)
	}
	bytes, e := ioutil.ReadAll(resp.Body)
	log.Info("pin res:", string(bytes), e)
	return e
}
