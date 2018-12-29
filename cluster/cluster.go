package cluster

import (
	"bufio"
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

// WaitingForInitialize ...
func WaitingForInitialize(ctx context.Context) bool {
	for {
		if !IsInitialized() {
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
		if initCheck(InitIPFS) {
			log.Println("init ipfs")
			firstRunIPFS()
		}
		if initCheck(InitService) {
			log.Println("init service")
			firstRunService()
		}
		StartIPFS(ctx)
		time.Sleep(5 * time.Second)
		StartService(ctx)
	}
}

// initCheck ...
func initCheck(name string) bool {
	file := cfg.RootPath + "/" + name
	info, err := os.Stat(file)
	log.Println(info)
	if err == nil {
		err := os.Remove(cfg.RootPath + "/" + name)
		if err == nil {
			return true
		}
	}
	return false
}

// StartIPFS ...
func StartIPFS(ctx context.Context) {
	go opRun(ctx, "ipfs", "daemon")
}

// StartService ...
func StartService(ctx context.Context) {
	if NeedBootstrap() {
		boot := getServiceBootstrap()
		if boot != "" {
			go opRun(ctx, cfg.ServiceCommandName, "daemon", "--bootstrap", boot)
			return
		}
	}
	go opRun(ctx, cfg.ServiceCommandName, "daemon")
}

// NeedBootstrap ...
func NeedBootstrap() bool {
	if cfg.HostType == HostClient {
		return true
	}
	return false
}

func getServiceBootstrap() string {
	url := strings.Join([]string{cfg.RemoteIP, "v0", "bootstrap"}, "/")
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

func opRun(ctx context.Context, command string, options ...string) error {
	cmd := exec.Command(command, options...)

	cmd.Env = cfg.GetEnv()

	//显示运行的命令
	log.Println(command, options)
	log.Println("command:", cmd.Args)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		errors.ErrorStack(err)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		errors.ErrorStack(err)
		return err
	}

	err = cmd.Start()
	if err != nil {
		errors.ErrorStack(err)
		return err
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))

	//实时循环读取输出流中的一行内容
	for {
		select {
		case <-ctx.Done():
			break
		}
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			break
		}

		log.Println(line)
	}

	err = cmd.Wait()
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	return err
}

// Reset ...
func Reset() error {

	for _, v := range cfg.ClusterEnviron {
		path := strings.Split(v, "=")[1]

		if strings.LastIndex(path, "/") == 0 {
			path = path + "/"
		}

		log.Println("clear", path)
		err := opRun(context.Background(), "rm", "-R", path)
		if err != nil {
			errors.ErrorStack(err)
			return err
		}
	}

	log.Println("clear /root/.ipfs-cluster-monitor")
	err := opRun(context.Background(), "rm", "-R", "/root/.ipfs-cluster-monitor")
	if err != nil {
		errors.ErrorStack(err)
		return err
	}

	isInitialized = false
	return nil
}
