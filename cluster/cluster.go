package cluster

import (
	"bufio"
	"context"
	"fmt"
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

// RunCMD ...
func RunCMD(command string, env []string, options ...string) error {
	cmd := exec.Command(command, options...)

	cmd.Env = env
	_, err := cmd.CombinedOutput()

	if err != nil {
		errors.ErrorStack(err)
		log.Println(err)
	}
	return err
}

func optimizeRunCMD(ctx context.Context, command string, env []string, options ...string) error {
	cmd := exec.CommandContext(ctx, command, options...)
	cmd.Env = env

	//显示运行的命令
	//log.Println("command:", cmd.Args)
	log.Output(2, fmt.Sprintln("command:", cmd.Args))

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return xerrors.Errorf("out pipe:%w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return xerrors.Errorf("err pipe:%w", err)
	}

	err = cmd.Start()
	if err != nil {
		return xerrors.Errorf("start:%w", err)
	}

	reader := bufio.NewReader(io.MultiReader(stdout, stderr))

	//实时循环读取输出流中的一行内容
	for {
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			log.Println("end", cmd.Args, e)
			break
		}

		log.Print(line)
	}

	err = cmd.Wait()
	if err != nil {
		return xerrors.Errorf("wait:%w", err)
	}
	return nil
}

func webAddress(api string) string {
	//url := strings.Join([]string{cfg.RemoteIP + cfg.RemotePort, cfg.Version, api}, "/")
	//return "http://" + url
	return ""
}

// InitMaker ...
func InitMaker(cfg *config.Configure) error {

	file, err := os.OpenFile(cfg.Root, os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if os.IsNotExist(err) {
		log.Println("not exist ", err)
		_ = os.MkdirAll(cfg.Root, os.ModePerm)
		file, err = os.OpenFile(cfg.FD(), os.O_RDWR|os.O_CREATE|os.O_SYNC, os.ModePerm)
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

	cfile, err := os.Create(filepath.Join(cfg.Root, config.Ipfs))
	log.Println("created:", cfile.Name())
	if err != nil {
		return xerrors.Errorf("ipfs file:%w", err)
	}
	defer cfile.Close()

	sfile, err := os.Create(filepath.Join(cfg.Root, config.Cluster))
	log.Println("created:", sfile.Name())
	if err != nil {
		return xerrors.Errorf("cluster file:%w", err)
	}
	defer sfile.Close()
	return nil
}

// InitRunning ...
func InitRunning(path string) bool {
	info, err := os.Stat(path)
	log.Printf("%+v,%+v\n", info, err) //has nil
	if err == nil {
		err := os.Remove(path)
		if err == nil {
			return true
		}
	}
	return false
}
