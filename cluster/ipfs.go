package cluster

import (
	"bufio"
	"fmt"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"github.com/juju/errors"
	"io"
	"log"
	"os"
	"os/exec"
)

func firstRunIPFS() {
	cmd := exec.Command(api.Config().CommandName, "init")
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

func optimizationFirstRunIPFS() {
	cmd := exec.Command(api.Config().CommandName, "init")
	cmd.Env = os.Environ()
	if clusterEnviron != nil {
		cmd.Env = append(cmd.Env, clusterEnviron...)
	}
	//显示运行的命令
	fmt.Println(cmd.Args)

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		errors.ErrorStack(err)
		return
	}

	err = cmd.Start()
	if err != nil {
		errors.ErrorStack(err)
		return
	}
	reader := bufio.NewReader(stdout)

	//实时循环读取输出流中的一行内容
	for {
		line, e := reader.ReadString('\n')
		if e != nil || io.EOF == e {
			break
		}
		fmt.Println(line)
	}

	err = cmd.Wait()
	if err != nil {
		errors.ErrorStack(err)
		return
	}
	return
}
