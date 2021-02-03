package main

import (
	"go-echo-server/handler"
	"go-echo-server/server"
)

func main() {
	go server.UdpListen(90, []handler.HandlerI{
		&handler.Logger{},
	})

	go server.HttpListen(80, []handler.HandlerI{
		&handler.Logger{},
	})

	for {
	}
}
