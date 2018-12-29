package cluster

import (
	"github.com/juju/errors"
	"io"
	"log"
	"os"
	"sync"
)

var once = sync.Once{}
var file *os.File

// LogInit ...
func LogInit() {
	log.Println("log initialized")
	var err error
	once.Do(func() {
		file, err = os.OpenFile(cfg.RootPath+"/monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
		if err == nil {
			output := io.MultiWriter(os.Stdout, file)
			log.SetOutput(output)
		} else {
			errors.ErrorStack(err)
		}
		log.SetFlags(log.Lshortfile | log.Ldate)
	})
}

// Log ...
func Log(v ...interface{}) {
	log.Println(v)
}

// ClearLog ...
func ClearLog() error {
	var err error
	file.Close()
	err = os.Remove(cfg.RootPath + "/monitor.log")
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	file, err = os.OpenFile(cfg.RootPath+"/monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
	if err != nil {
		output := io.MultiWriter(os.Stdout, file)
		log.SetOutput(output)
	}
	log.SetFlags(log.Lshortfile | log.Ldate)
	return nil
}
