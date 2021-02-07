package server

import (
	"go-echo-server/handler"
	"net"
)

func process(buf []byte, addr *net.UDPAddr, handlers []handler.Handler) {
	projectName := ""
	content := buf

	if buf[0] == '$' {
		for i := 2; i < 32; i++ {
			if buf[i] == '#' {
				projectName = string(buf[1:i])
				content = buf[i+1:]
			}
		}
	}

	for _, handler := range handlers {
		handler.Handle("udp", addr.String(), string(content), projectName)
	}
}

// UDPListen *
func UDPListen(port int, h []handler.Handler) {
	serverConn, _ := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	defer serverConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, _ := serverConn.ReadFromUDP(buf)
		go process(buf[0:n], addr, h)
	}
}
