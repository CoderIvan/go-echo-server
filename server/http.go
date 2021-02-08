package server

import (
	"go-echo-server/handler"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTPServer *
type HTTPServer struct {
	Port int
}

// Listen *
func (server *HTTPServer) Listen(h handler.Handler) {
	r := gin.New()

	f := func(c *gin.Context) {
		projectName := c.Param("projectName")

		buf := make([]byte, 1024)
		n, _ := c.Request.Body.Read(buf)

		h.Handle("http", c.Request.RemoteAddr, string(buf[0:n]), projectName)
	}

	r.POST("/", f)
	r.POST("/:projectName", f)

	r.Run(":" + strconv.Itoa(server.Port))
}
