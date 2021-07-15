package server

import (
	"go-echo-server/datagram"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// HTTPServer *
type httpServer struct {
	port   int
	server *http.Server
}

func NewHTTPServer(port int) *httpServer {
	return &httpServer{
		port: port,
	}
}

func setupRouter(handle func(datagram.Datagram)) *gin.Engine {
	router := gin.New()

	f := func(c *gin.Context) {
		projectName := c.Param("projectName")

		buf := make([]byte, 1024)
		n, _ := c.Request.Body.Read(buf)

		handle(datagram.Datagram{
			TagName:     "http-server",
			Addr:        c.Request.RemoteAddr,
			ProjectName: projectName,
			Content:     string(buf[0:n]),
			Time:        time.Now().UnixNano(),
		})
	}

	router.POST("/", f)
	router.POST("/:projectName", f)

	return router
}

// Listen *
func (hServer *httpServer) Listen(handle func(datagram.Datagram)) error {
	router := setupRouter(handle)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(hServer.port),
		Handler: router,
	}
	return server.ListenAndServe()
}

func (hServer *httpServer) Close() {
	if hServer.server != nil {
		hServer.server.Close()
	}
}
