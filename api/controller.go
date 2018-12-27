package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitPost ...
func InitPost(s string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		remote := ctx.PostForm("remote")
		secret := ctx.PostForm("secret")
		if !isInitialized {
			if remote != "" {
				config.SetClient(remote)
				config.Secret = secret
			} else {
				config.Secret = GenerateRandomString(32)
			}

			config.InitLoader()
		}

		ctx.JSON(http.StatusOK, config)
	}
}

// HeartBeatGet ...
func HeartBeatGet(s string) gin.HandlerFunc {

	return func(ctx *gin.Context) {

	}
}
