package api

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/cluster"
)

// Router ...
func Router(eng *gin.Engine) {

	ver := cluster.Config().Version

	v0 := eng.Group(ver)

	v0.POST("init", InitPost(ver))

	v0.GET("heartbeat", HeartBeatGet(ver))

	v0.GET("log", LogGet(ver))

	v0.GET("bootstrap", BootstrapGet(ver))

	v0.GET("reset", ResetGet(ver))

	v0.GET("delete/:id", DeleteGet(ver))
}
