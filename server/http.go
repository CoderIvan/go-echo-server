package server

import (
	"go-echo-server/handler"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTPListen *
func HTTPListen(port int, handlers []handler.Handler) {
	r := gin.New()

	f := func(c *gin.Context) {
		projectName := c.Param("projectName")

		buf := make([]byte, 1024)
		n, _ := c.Request.Body.Read(buf)

		for _, handler := range handlers {
			handler.Handle("http", c.Request.RemoteAddr, string(buf[0:n]), projectName)
		}
	}

	r.POST("/", f)
	r.Run(":" + strconv.Itoa(port))
}
