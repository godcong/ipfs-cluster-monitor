package cluster

import (
	"bufio"
	"context"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"github.com/pelletier/go-toml"
	"golang.org/x/exp/xerrors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StatusCode ...
type StatusCode int

// StatusFailed ...
const (
	StatusFailed      StatusCode = -1
	StatusSuccess     StatusCode = 0
	StatusStart       StatusCode = iota
	StatusCreated     StatusCode = iota
	StatusProcessing  StatusCode = iota
	StatusIpfsInit    StatusCode = iota
	StatusServiceInit StatusCode = iota
	StatusIpfsRun     StatusCode = iota
	StatusServiceRun  StatusCode = iota
)

// ResultMessage ...
type ResultMessage struct {
	Code    int                    `json:"code"`
	Detail  map[string]interface{} `json:"detail"`
	Message string                 `json:"message"`
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

	//cmd.Env = append(cmd.Env)
	_, err := cmd.CombinedOutput()

	if err != nil {
		errors.ErrorStack(err)
		log.Println(err)
	}
	return err
}

func optimizeRunCMD(ctx context.Context, command string, options ...string) error {
	cmd := exec.CommandContext(ctx, command, options...)
	//end := strconv.FormatInt(time.Now().Unix(), 10)
	//c.commands.Store(command+"_"+end, cmd)

	//cmd.Env = cfg.GetEnv()

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
	//c.commands.Range(
	//	func(key, value interface{}) bool {
	//		if v, b := value.(*exec.Cmd); b {
	//			log.Println("kill", key)
	//			err := v.Process.Kill()
	//			if err != nil {
	//				errors.ErrorStack(err)
	//				log.Println(err)
	//			}
	//			c.commands.Delete(key)
	//			return true
	//		}
	//		log.Println(key, "not cmd continue")
	//		return true
	//	})
}

func webAddress(api string) string {
	//url := strings.Join([]string{cfg.RemoteIP + cfg.RemotePort, cfg.Version, api}, "/")
	//return "http://" + url
	return ""
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
	//c.Stop()

	clear(config.IpfsPath())
	clear(config.IpfsClusterPath())
	//clear(cfg.RootPath)

	//reset config
	//cfg = DefaultConfig()

	//reset status
	//c.isInitialized = false
	//c.SetStatus("init", StatusFailed)

	//waiting 30 sec to restart
	//for ; waiting >= 0; waiting-- {
	//	atomic.StoreInt32(&c.waiting, waiting)
	//	time.Sleep(time.Second)
	//}
	//
	//rerun
	//go c.Start()
	return nil
}

// InitMaker ...
func InitMaker(cfg *config.Configure, configPath string) error {
	dir, _ := filepath.Split(configPath)
	file, err := os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if os.IsNotExist(err) {
		log.Println("not exist ", err)
		_ = os.MkdirAll(dir, os.ModePerm)
		file, err = os.OpenFile(configPath, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
		if err != nil {
			return xerrors.Errorf("make file:%w", err)
		}
	} else {

	}

	defer file.Close()
	enc := toml.NewEncoder(file)
	err = enc.Encode(*cfg)
	log.Println("created:", file.Name())
	if err != nil {
		return xerrors.Errorf("encode file:%w", err)
	}

	cfile, err := os.Create(filepath.Join(dir, config.InitIPFS))
	log.Println("created:", cfile.Name())
	if err != nil {
		return xerrors.Errorf("ipfs file:%w", err)
	}
	defer cfile.Close()

	sfile, err := os.Create(filepath.Join(dir, config.InitIPFSCluster))
	log.Println("created:", sfile.Name())
	if err != nil {
		return xerrors.Errorf("cluster file:%w", err)
	}
	defer sfile.Close()
	return nil
}
