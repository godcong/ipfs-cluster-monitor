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
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				findMyFather()
				time.Sleep(cfg.ServerCheckInterval)
			}
			//get info
		}
	}()
}

func findMyFather() error {
	q := url.Values{}
	q.Set("name", "son")
	response, err := http.PostForm(webAddress("join"), q)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	log.Println(string(bytes))
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	var rm ResultMessage

	err = jsoniter.Unmarshal(bytes, &rm)
	if err != nil {
		errors.ErrorStack(err)
		log.Println(err.Error())
		return err
	}
	if rm.Code == 0 {
		log.Println("join success")
		return nil
	}
	log.Println("failed:", rm.Message)
	return errors.New(rm.Message)
}

// AddMySon ...
func AddMySon(key, val string) {
	if _, loaded := clients.LoadOrStore(key, val); loaded {
		log.Println(key, "is already joined")
	} else {
		log.Println(key, "is joined")
	}
}
