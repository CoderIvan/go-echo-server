package main

import (
	"go-echo-server/handler"
	"go-echo-server/server"
)

type connect struct {
	servers  []server.Server
	handlers []handler.Handler
}

func (ct connect) Handle(protocol string, addr string, content string, projectName string) {
	for _, handler := range ct.handlers {
		handler.Handle(protocol, addr, content, projectName)
	}
}

func (ct connect) run() {
	for _, server := range ct.servers {
		go server.Listen(ct)
	}
}

func main() {
	ct := connect{
		[]server.Server{
			&server.UDPServer{
				Port: 90,
			},
		},
		[]handler.Handler{
			&handler.Logger{},
			&handler.Sls{},
		},
	}

	ct.run()

	// go udpServer.Listen(handler)

	// go server.HTTPListen(80, []handler.Handler{
	// 	&handler.Logger{},
	// 	&handler.Sls{},
	// })

	for {
	}
}
