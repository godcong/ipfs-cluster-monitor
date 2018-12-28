package api

import (
	"io"
	"log"
	"os"
	"sync"
)

var once = sync.Once{}

// Log ...
func Log(v ...interface{}) {
	once.Do(func() {
		file, err := os.OpenFile(cfg.RootPath+"/monitor.log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
		if err != nil {
			output := io.MultiWriter(os.Stdout, file)
			log.SetOutput(output)
		}
		log.SetFlags(log.Lshortfile | log.Ldate)
	})

	log.Println(v)

}
