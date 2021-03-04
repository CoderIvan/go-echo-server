package server

import (
	"go-echo-server/datagram"
	"net"
	"time"
)

// UDPServer *
type UDPServer struct {
	Port int
}

func process(buf []byte, addr *net.UDPAddr) datagram.Datagram {
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

	return datagram.Datagram{
		TagName:     "udp-server",
		Addr:        addr.String(),
		ProjectName: projectName,
		Content:     string(content),
		Time:        time.Now().Unix(),
	}
}

// Listen *
func (server *UDPServer) Listen(ch chan datagram.Datagram) {
	serverConn, _ := net.ListenUDP("udp", &net.UDPAddr{
		Port: server.Port,
	})
	defer serverConn.Close()

	for {
		buf := make([]byte, 1024)
		n, addr, _ := serverConn.ReadFromUDP(buf)

		ch <- process(buf[0:n], addr)
	}
}
