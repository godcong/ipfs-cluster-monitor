package main

import (
	"github.com/gin-gonic/gin"
	"github.com/godcong/ipfs-cluster-monitor/api"
)

func main() {
	engine := gin.Default()
	api.DefaultConfig()
	api.Router(engine)

}
