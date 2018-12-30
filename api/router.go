package api

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
	"io/ioutil"
	"net/http"
)

// Router ...
func Router(eng *gin.Engine) {

	ver := cluster.Config().Version

	v0 := eng.Group(ver)

	v0.POST("init", InitPost(ver))
	v0.GET("init", InitGet(ver))
	v0.GET("heartbeat", HeartBeatGet(ver))

	v0.GET("log", LogGet(ver))

	v0.GET("reset", ResetGet(ver))
	v0.GET("waiting", WaitingGet(ver))
	v0.GET("delete/:pn", DeleteGet(ver))
	v0.POST("join", JoinPost(ver))

	v0.GET("bootstrap", BootstrapGet(ver))
	v0.GET("secret", SecretGet(ver))
	v0.GET("killyou/:id", KillGet(ver))

	v0.Any("debug", func(context *gin.Context) {
		request, err := http.NewRequest(context.Request.Method, context.Query("url"), context.Request.Body)
		if err != nil {
			failed(context, err.Error())
			return
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			failed(context, err.Error())
			return
		}
		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			failed(context, err.Error())
			return
		}
		context.String(http.StatusOK, string(bytes))
	})
}
