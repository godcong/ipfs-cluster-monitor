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
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ResultMessage ...
type ResultMessage struct {
	Code    int                    `json:"code"`
	Detail  map[string]interface{} `json:"detail"`
	Message string                 `json:"message"`
}

type Cluster struct {
	GlobalContext context.Context
	context       context.Context
	cancelFunc    context.CancelFunc
	commands      sync.Map
	status        sync.Map
	isInitialized bool
	waiting       int32
}

var cluster *Cluster
var globalContext context.Context

func init() {
	cluster = newCluster()
}

func newCluster() *Cluster {
	return &Cluster{}
}

func DefaultCluster() *Cluster {
	return cluster
}

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
	cluster.GlobalContext = ctx
	cluster.context, cluster.cancelFunc = context.WithCancel(context.Background())

	if waitingForInitialize(cluster.context) {
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
		runIPFS(globalContext)
		//var service context.Context
		//service, cancelService = context.WithCancel(context.Background())
		runService(globalContext)

		time.Sleep(cfg.Interval)

		if isClient() {
			runJoin(cluster.context)
		} else {
			runMonitor(cluster.context)
		}
		atomic.StoreInt32(&cluster.waiting, -1)
	}

}

// initCheck ...
func initCheck(name string) bool {
	path := filepath.Join(cfg.RootPath, name)
	info, err := os.Stat(path)
	log.Println(info.Name())
	if err == nil {
		err := os.Remove(path)
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

	if err != nil {
		errors.ErrorStack(err)
		log.Println(err)
	}
	return err
}

func (c *Cluster) optimizeRunCMD(command string, options ...string) error {
	cmd := exec.Command(command, options...)
	end := strconv.FormatInt(time.Now().Unix(), 10)
	c.commands.Store(command+"_"+end, cmd)

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
func (c *Cluster) stopRunningCMD() {
	c.commands.Range(
		func(key, value interface{}) bool {
			if v, b := value.(*exec.Cmd); b {
				log.Println("kill", key)
				err := v.Process.Kill()
				if err != nil {
					errors.ErrorStack(err)
					log.Println(err)
				}
				c.commands.Delete(key)
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

func (c *Cluster) ResetWaiting() int {
	return int(atomic.LoadInt32(&c.waiting))
}

// Reset ...
func (c *Cluster) Reset() error {
	waiting := int32(cfg.ResetWaiting)
	//stop running ipfs and service
	c.stopRunningCMD()

	clear(ipfsPath())
	clear(servicePath())
	clear(cfg.RootPath)

	//reset config
	cfg = DefaultConfig()

	if c.cancelFunc != nil {
		c.cancelFunc()
		c.cancelFunc = nil
	}

	//reset status

	c.isInitialized = false
	SetStatus("init", StatusFailed)

	//waiting 30 sec to restart
	for ; waiting >= 0; waiting-- {
		atomic.StoreInt32(&c.waiting, waiting)
		time.Sleep(time.Second)
	}

	//rerun
	go Run(nil)
	return nil
}

type StatusCode int

const (
	StatusFailed     StatusCode = -1
	StautsSuccess    StatusCode = 0
	StatusProcessing StatusCode = 1
)

func SetStatus(key string, value StatusCode) {
	status.Store(key, value)
}

func GetStatus(key string) StatusCode {
	value, ok := status.Load(key)
	if !ok {
		return StatusFailed
	}
	return value.(StatusCode)
}
