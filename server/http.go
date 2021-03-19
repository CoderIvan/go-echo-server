package server

import (
	"go-echo-server/datagram"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPServer *
type HTTPServer struct {
	Port int
}

// Listen *
func (server *HTTPServer) Listen(ch chan datagram.Datagram) {
	r := gin.New()

	f := func(c *gin.Context) {
		projectName := c.Param("projectName")

		buf := make([]byte, 1024)
		n, _ := c.Request.Body.Read(buf)

		ch <- datagram.Datagram{
			TagName:     "http-server",
			Addr:        c.Request.RemoteAddr,
			ProjectName: projectName,
			Content:     string(buf[0:n]),
			Time:        time.Now().Unix(),
		}
	}

	r.POST("/", f)
	r.POST("/:projectName", f)

	r.Run(":" + strconv.Itoa(server.Port))
}
