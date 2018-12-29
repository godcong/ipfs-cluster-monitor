package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

//runJoin join remote server
func runJoin(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			joinToServer()
			time.Sleep(cfg.ServerCheckInterval)
		}
		//get info
	}
}

func joinToServer() int {
	response, err := http.PostForm(webAddress("join"), url.Values{})
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
	var rm ResultMessage

	err = jsoniter.Unmarshal(bytes, rm)
	if err != nil {
		errors.ErrorStack(err)
		panic(err)
	}
	if rm.Code == 0 {
		log.Println("join success")
	}
	log.Println("failed:", rm.Message)
	return rm.Code
}

// JoinFromClient ...
func JoinFromClient(key, val string) {
	if _, loaded := clients.LoadOrStore(key, val); loaded {
		log.Println(key, "is already joined")
	} else {
		log.Println(key, "is joined")
	}
}
