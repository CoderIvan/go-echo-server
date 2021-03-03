package server

import (
	"go-echo-server/datagram"
	"go-echo-server/handler"
	"net"
)

// UDPServer *
type UDPServer struct {
	Port int
}

func process(buf []byte, addr *net.UDPAddr, h handler.Handler) {
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

	h.Handle(datagram.Datagram{
		TagName:     "udp-server",
		Addr:        addr.String(),
		ProjectName: projectName,
		Content:     string(content),
	})
}

// Listen *
func (server *UDPServer) Listen(h handler.Handler) {
	serverConn, _ := net.ListenUDP("udp", &net.UDPAddr{
		Port: server.Port,
	})
	defer serverConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, _ := serverConn.ReadFromUDP(buf)
		go process(buf[0:n], addr, h)
	}
}
