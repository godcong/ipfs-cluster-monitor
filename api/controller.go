package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// InitGet ...
func InitGet(ver string) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	}

}
