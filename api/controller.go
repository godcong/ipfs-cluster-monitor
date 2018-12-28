package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

// InitPost ...
func InitPost(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote := ctx.PostForm("Remote")
		secret := ctx.PostForm("Secret")
		ipfs := ctx.PostForm("IPFS_PATH")
		service := ctx.PostForm("IPFS_CLUSTER_PATH")
		//monitor := ctx.PostForm("IPFS_CLUSTER_MONITOR")
		if !IsInitialized() {
			if remote != "" {
				cfg.SetClient(remote)
				cfg.Secret = secret
			} else {
				cfg.Secret = GenerateRandomString(64)
			}
			if ipfs != "" {
				cfg.MonitorEnviron = append(cfg.MonitorEnviron, ipfs)
			}
			if service != "" {
				cfg.MonitorEnviron = append(cfg.MonitorEnviron, service)
			}
			//if monitor != "" {
			//	cfg.MonitorEnviron = append(cfg.MonitorEnviron, monitor)
			//}

			cfg.Make()
			log.Println("host initialized")

			ctx.JSON(http.StatusOK, cfg)
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "host is initialized"})
	}
}

// HeartBeatGet ...
func HeartBeatGet(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		response, err := http.Get("http://localhost:9094/peers")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		bytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.String(http.StatusOK, string(bytes))

	}
}
