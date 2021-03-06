package cluster

import (
	"context"
	"github.com/json-iterator/go"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
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
				//time.Sleep(cfg.ServerCheckInterval)
			}
			//get info
		}
	}()
}

func findMyFather() error {
	q := url.Values{}

	service, err := getServiceInfo()
	if err != nil {
		return err
	}
	q.Set("name", service.Peername)
	response, err := http.PostForm(webAddress("join"), q)
	if err != nil {
		errors.ErrorStack(err)
		return err
	}
	bytes, err := ioutil.ReadAll(response.Body)
	log.Println("join", string(bytes))
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
		return nil
	}
	return errors.New(rm.Message)
}
