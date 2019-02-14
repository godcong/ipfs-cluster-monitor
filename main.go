//go:generate protoc --go_out=plugins=grpc:./proto monitor.proto
package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/service"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var configPath = flag.String("config", "config.toml", "config path")
var logPath = flag.String("log", "monitor.log", "log path")

func main() {

	flag.Parse()

	dir, _ := filepath.Split(*logPath)
	_ = os.MkdirAll(dir, os.ModePerm)

	file, err := os.OpenFile(*logPath, os.O_SYNC|os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	err = config.Initialize(*configPath)
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
