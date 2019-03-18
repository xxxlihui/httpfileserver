package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var UserName string
var Password string

func author(c *gin.Context)  {
	username,password,ok:=c.Request.BasicAuth()
	if !ok||username!=UserName||password!=Password {
		c.Header("WWW-Authenticate", `Basic realm="Dotcoo User Login"`)
		c.AbortWithStatus(http.StatusUnauthorized)
	}else {
		c.Next()
	}
}
