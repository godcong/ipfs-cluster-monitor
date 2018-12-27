package main

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/api"
	"github.com/juju/errors"
	"net/http"
)

func main() {
	engine := gin.Default()
	api.DefaultConfig().InitLoader()
	engine.NoRoute(NoResponse)
	api.Router(engine)
	err := engine.Run()
	if err != nil {
		errors.ErrorStack(err)
	}
}

// NoResponse ...
func NoResponse(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{
		"code":    -1,
		"message": "remote address not found",
	})
}
