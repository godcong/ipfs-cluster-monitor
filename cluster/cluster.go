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
	"sync"
	"time"
)

var cluster sync.Map
var globalContext context.Context

// waitingForInitialize ...
func waitingForInitialize(ctx context.Context) bool {
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
	globalContext = ctx

	if waitingForInitialize(ctx) {
		if initCheck(InitIPFS) {
			log.Println("init ipfs")
			firstRunIPFS()
		}
		if initCheck(InitService) {
			log.Println("init service")
			firstRunService()
		}
		//var ipfs context.Context
		//ipfs, cancelIPFS = context.WithCancel(context.Background())
		runIPFS(ctx)
		time.Sleep(5 * time.Second)
		//var service context.Context
		//service, cancelService = context.WithCancel(context.Background())
		runService(ctx)
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

// useBootstrap ...
func useBootstrap() bool {
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

func runCMD(command string, options ...string) error {
	cmd := exec.Command(command, options...)

	cmd.Env = cfg.GetEnv()
	_, err := cmd.CombinedOutput()
	//if bts != nil {
	//	bts = bytes.TrimSpace(bts)
	//	log.Println(string(bts))
	//}

	if err != nil {
		errors.ErrorStack(err)
		log.Println(err)
	}
	return err
}

func optimizeRunCMD(command string, options ...string) error {
	cmd := exec.Command(command, options...)
	cluster.Store(command, cmd)

	cmd.Env = cfg.GetEnv()

	//显示运行的命令
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
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			log.Println("end", cmd.Args, e)
			errors.ErrorStack(err)
			break
		}

		log.Print(line)
	}

	err = cmd.Wait()
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	return err
}

// stopRunningCMD ...
func stopRunningCMD() {
	cluster.Range(
		func(key, value interface{}) bool {
			if v, b := value.(*exec.Cmd); b {
				log.Println("kill", key)
				err := v.Process.Kill()
				if err != nil {
					errors.ErrorStack(err)
					log.Println(err)
					return true
				}
				cluster.Delete(key)
				return true
			}
			log.Println(key, "not cmd continue")
			return true
		})
}

// Reset ...
func Reset() error {
	stopRunningCMD()
	for _, v := range cfg.ClusterEnviron {
		path := strings.Split(v, "=")[1]

		if strings.LastIndex(path, "/") != 0 {
			path = path + "/"
		}

		log.Println("clear", path)
		err := runCMD("rm", "-R", path)
		if err != nil {
			errors.ErrorStack(err)
			continue
		}
	}

	log.Println("clear", cfg.RootPath+"/")
	err := runCMD("rm", "-R", cfg.RootPath+"/")
	if err != nil {
		errors.ErrorStack(err)
	}
	cfg = DefaultConfig()
	isInitialized = false
	go Run(globalContext)
	return nil
}
