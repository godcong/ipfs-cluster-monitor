package api

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const prefix = "ICM"

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
				cfg.Secret = prefix + GenerateRandomString(64)
			}
			cfg.SetEnv("IPFS_PATH", ipfs)
			cfg.SetEnv("IPFS_CLUSTER_PATH", service)
			cfg.Make()

			log.Println("host initialized")

			ctx.JSON(http.StatusOK, cfg)
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "host is initialized"})
	}
}

// HeartBeatGet ...
func HeartBeatGet(ver string) gin.HandlerFunc {
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

// LogGet ...
func LogGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

// BootstrapGet ...
func BootstrapGet(ver string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error
		var s ServiceStatus
		resp, err := http.Get("http://127.0.0.1:9094/id")
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		bytes, err := ioutil.ReadAll(resp.Body)
		err = jsoniter.Unmarshal(bytes, &s)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		address := ""

		for _, value := range s.Addresses {
			if strings.Index(value, "/p2p-circuit") >= 0 {
				continue
			} else if strings.Index(value, "/ip4/127.0.0.1") >= 0 {
				continue
			}
			address = value
			break
		}
		ctx.JSON(http.StatusOK, gin.H{"bootstrap": address})
	}
}
