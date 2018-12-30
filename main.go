package main

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"github.com/juju/errors"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	cluster.LogInit()

	engine := gin.Default()
	engine.NoRoute(NoResponse)

	c := cluster.Default()
	go c.Start()
	defer c.Stop()

	api.Router(engine)
	err := engine.Run(cluster.Config().RemotePort)
	if err != nil {
		errors.ErrorStack(err)
	}
}

// NoResponse ...
func NoResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    -1,
		"message": "remote address not found",
	})
}

func testCommand() {
	//D:\\workspace\\project\\docker\\
	cmd := exec.Command("nohup", "ipfs", "cluster", "&")

	cmd.Env = os.Environ()

	err := cmd.Run()
	println(err)
	log.Println(cmd.Output())
	//bytes, err := cmd.CombinedOutput()
	//log.Println(bytes, err)

}
