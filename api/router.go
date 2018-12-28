package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Router ...
func Router(eng *gin.Engine) {

	ver := cfg.Version

	v0 := eng.Group(ver)

	v0.POST("init", InitPost(ver))

	v0.GET("heartbeat", HeartBeatGet(ver))

	v0.GET("log", LogGet(ver))

	v0.GET("bootstrap", BootstrapGet(ver))

	v0.GET("reset", ResetGet(ver))
}

func ResetGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := gin.H{}

		ctx.JSON(http.StatusOK)
	}
}
