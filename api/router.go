package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

// Router ...
func Router(eng *gin.Engine) {

	ver := "v0"

	v0 := eng.Group(ver)

	v0.POST("init", InitPost(ver))

	v0.GET("heartbeat", HeartBeatGet(ver))

	v0.GET("log", LogGet(ver))

}
