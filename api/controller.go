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
		remote := ctx.PostForm("remote")
		secret := ctx.PostForm("secret")
		if !IsInitialized(config) {
			if remote != "" {
				config.SetClient(remote)
				config.Secret = secret
			} else {
				config.Secret = GenerateRandomString(32)
			}

			config.Make()

			log.Println("initialized")
		}

		ctx.JSON(http.StatusOK, config)
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
