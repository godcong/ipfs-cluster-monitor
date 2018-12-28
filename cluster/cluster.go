package cluster

import (
	"context"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var clusterEnviron []string

func firstRunIPFS() {
	cmd := exec.Command("ipfs", "init")
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

// Init ...
func Init() {
	clusterEnviron = api.Config().ClusterEnviron
}

// Run ...
func Run(ctx context.Context) {
	if WaitingForInitialize(ctx) {
		Init()
		if NeedInit(api.InitIPFS) {
			firstRunIPFS()
		}
		if NeedInit(api.InitService) {
			firstRunService()
		}
		StartIPFS(ctx)
		time.Sleep(5 * time.Second)
		StartService(ctx)
	}
}

// NeedInit ...
func NeedInit(name string) bool {
	file := api.Config().RootPath + "/" + name
	info, err := os.Stat(file)
	log.Println(info)
	if err == nil {
		err := os.Remove(api.Config().RootPath + "/" + name)
		if err == nil {
			return true
		}
	}
	return false
}

// StartIPFS ...
func StartIPFS(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := exec.Command("ipfs", "daemon")
			cmd.Env = os.Environ()

			err := cmd.Run()
			if err != nil {
				errors.ErrorStack(err)
				log.Println(err)
				return
			}
		}
	}()

}

// StartService ...
func StartService(ctx context.Context) {
	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := exec.Command("ipfs-cluster-service", "daemon")
			if NeedBootstrap() {
				boot := getServiceBootstrap()
				if boot != "" {
					cmd = exec.Command("ipfs-cluster-service", "daemon", "--bootstrap", boot)
				}
			}

			cmd.Env = os.Environ()
			err := cmd.Run()
			if err != nil {
				errors.ErrorStack(err)
				log.Println(err)
				return
			}
		}
	}()
}

// NeedBootstrap ...
func NeedBootstrap() bool {
	log.Println(api.Config().HostType)
	if api.Config().HostType == api.HostClient {
		return true
	}
	return false
}

func getServiceBootstrap() string {
	url := strings.Join([]string{api.Config().RemoteIP, "v0", "bootstrap"}, "/")
	response, err := http.Get("http://" + url)
	if err != nil {
		return ""
	}
	bytes, err := ioutil.ReadAll(response.Body)
	var m map[string]string
	err = jsoniter.Unmarshal(bytes, &m)
	if err != nil {
		return ""
	}
	v, b := m["bootstrap"]
	if b {
		return v
	}
	return ""
}
