package main

import (
	"go-echo-server/handler"
	"go-echo-server/server"
)

func main() {
	go server.UDPListen(90, []handler.Handler{
		&handler.Logger{},
		&handler.Sls{},
	})

	// go server.HTTPListen(80, []handler.Handler{
	// 	&handler.Logger{},
	// 	&handler.Sls{},
	// })

	for {
	}
}
