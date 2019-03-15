//go:generate protoc --go_out=plugins=grpc:./proto monitor.proto
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/godcong/go-trait"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var configPath = flag.String("config", "config.toml", "config path")
var logPath = flag.String("log", "logs/monitor.log", "log path")
var level = flag.String("level", "debug", "set log output level")

func main() {

	flag.Parse()

	trait.InitRotateLog(*logPath)

	err := config.Initialize(os.Args[0], *configPath)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	//start
	service.Start()

	go func() {
		sig := <-sigs
		fmt.Println(sig, "exiting")
		service.Stop()
		done <- true
	}()
	<-done
}

// NoResponse ...
func NoResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    -1,
		"message": "remote address not found",
	})
}
