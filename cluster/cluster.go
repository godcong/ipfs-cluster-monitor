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

// ResultMessage ...
type ResultMessage struct {
	Code    int                    `json:"code"`
	Detail  map[string]interface{} `json:"detail"`
	Message string                 `json:"message"`
}

var commands sync.Map
var globalContext context.Context
var globalCancel context.CancelFunc

// waitingForInitialize ...
func waitingForInitialize(ctx context.Context) bool {
	for {
		if !IsInitialized() {
			time.Sleep(cfg.Interval)
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
	globalContext, globalCancel = context.WithCancel(context.Background())

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
		//var service context.Context
		//service, cancelService = context.WithCancel(context.Background())
		runService(ctx)

		time.Sleep(cfg.Interval)
		if isClient() {
			runJoin(globalContext)
		} else {
			runMonitor(globalContext)
		}

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

// isClient ...
func isClient() bool {
	if cfg.HostType == HostClient {
		return true
	}
	return false
}

func getServiceBootstrap() string {

	response, err := http.Get(webAddress("bootstrap"))
	if err != nil {
		return ""
	}
	bytes, err := ioutil.ReadAll(response.Body)
	var msg ResultMessage
	err = jsoniter.Unmarshal(bytes, &msg)
	if err != nil {
		return ""
	}
	v, b := msg.Detail["bootstrap"]
	if b {
		if v1, b := v.(string); b {
			return v1
		}
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
	commands.Store(command, cmd)

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
	commands.Range(
		func(key, value interface{}) bool {
			if v, b := value.(*exec.Cmd); b {
				log.Println("kill", key)
				err := v.Process.Kill()
				if err != nil {
					errors.ErrorStack(err)
					log.Println(err)
					return true
				}
				commands.Delete(key)
				return true
			}
			log.Println(key, "not cmd continue")
			return true
		})
}

func webAddress(api string) string {
	url := strings.Join([]string{cfg.RemoteIP + cfg.RemotePort, cfg.Version, api}, "/")
	return "http://" + url
}

func clear(path string) {
	if strings.LastIndex(path, "/") != 0 {
		path = path + "/"
	}
	log.Println("clear", path)
	err := runCMD("rm", "-R", path)
	if err != nil {
		errors.ErrorStack(err)
	}
	return
}

// Reset ...
func Reset() error {
	//stop running ipfs and service
	stopRunningCMD()

	clear(ipfsPath())
	clear(servicePath())
	clear(cfg.RootPath)

	//reset config
	cfg = DefaultConfig()
	//reset status
	isInitialized = false

	if globalCancel != nil {
		globalCancel()
		globalCancel = nil
	}

	//rerun
	go Run(globalContext)
	return nil
}
