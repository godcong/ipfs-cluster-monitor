package api

import "github.com/gin-gonic/gin"

// Router ...
func Router(eng *gin.Engine) {

	ver := "v0"

	v0 := eng.Group(ver)

	v0.GET("init", InitGet(ver))

}
