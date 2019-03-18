package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func fileSecurity(c *gin.Context) {
	if strings.HasPrefix(c.Request.RequestURI, "/..") {
		c.String(http.StatusOK, "非法访问")
		c.Abort()
	}
}
