package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

// IpfsInfo ...
type IpfsInfo struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

// runIPFS ...
func runIPFS(ctx context.Context) {
	go optimizeRunCMD(ctx, "ipfs", "daemon")
}

func firstRunIPFS() {
	//cmd := exec.Command(cfg.CommandName, "init")
	//cmd.Env = cfg.GetEnv()

	//bytes, err := cmd.CombinedOutput()
	//log.Println(string(bytes))
	//if err != nil {
	//	errors.ErrorStack(err)
	//	panic(err)
	//}
	return
}

func optimizationFirstRunIPFS(ctx context.Context) {
	err := optimizeRunCMD(ctx, "ipfs", "init")
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

func waitingIpfs(ctx context.Context) {
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:

		}
		_, err = getIpfsInfo()
		if err == nil {
			break
		}
		//time.Sleep(cfg.Interval)
	}
}
