//go:generate protoc --go_out=plugins=grpc:./proto monitor.proto
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/config"
	"github.com/godcong/ipfs-cluster-monitor/service"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var configPath = flag.String("config", "config.toml", "config path")
var logPath = flag.String("log", "logs/monitor.log", "log path")
var level = flag.String("level", "debug", "set log output level")

func main() {

	flag.Parse()

	InitLogger(*logPath, *level)

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

// InitLogger ...
func InitLogger(logPath string, level string) {
	dir, fname := filepath.Split(logPath)
	_ = os.MkdirAll(dir, os.ModePerm)
	writer, err := rotatelogs.New(
		dir+"%Y%m%d%H%M_"+fname,
		rotatelogs.WithLinkName(logPath),           // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),      // 文件最大保存时间
		rotatelogs.WithRotationTime(3*time.Minute), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %v", errors.WithStack(err))
	}

	//log.SetFormatter(&log.TextFormatter{})
	switch level {
	/*
		如果日志级别不是debug就不要打印日志到控制台了
	*/
	case "debug":
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stderr)
	case "info":
		setNull()
		log.SetLevel(log.InfoLevel)
	case "warn":
		setNull()
		log.SetLevel(log.WarnLevel)
	case "error":
		setNull()
		log.SetLevel(log.ErrorLevel)
	default:
		setNull()
		log.SetLevel(log.InfoLevel)
	}

	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, &log.JSONFormatter{})
	log.AddHook(lfHook)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
}

func setNull() {
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}
	writer := bufio.NewWriter(src)
	log.SetOutput(writer)
}
