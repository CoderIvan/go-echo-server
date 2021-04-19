package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-echo-server/datagram"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ParserFactory(callback func([]byte)) func(b []byte) {
	var bufferBuilder []byte

	return func(b []byte) {
		bufferBuilder = append(bufferBuilder, b...)

		var circle func()
		circle = func() {
			length := uint16(binary.BigEndian.Uint16(bufferBuilder[0:2]))
			if len(bufferBuilder) >= 2+int(length) {
				content := bufferBuilder[2 : length+2]
				bufferBuilder = bufferBuilder[length+2:]
				callback(content)
				if len(bufferBuilder) > 0 {
					circle()
				}
			}
		}
		circle()
	}
}

// UDPServer *
type tcpServer struct {
	listener net.Listener
}

func NewTCPServer(port int) *tcpServer {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		panic(fmt.Errorf("Listen tcp server failed, %v", err))
	}
	return &tcpServer{
		listener: listener,
	}
}

func handleConn(conn net.Conn, handle func(datagram.Datagram)) {
	defer conn.Close()

	var fistPkg bool = true
	addr := conn.RemoteAddr().String()
	var projectName string
	var contextID string
	connectID := uuid.New().String()

	extraInfo, _ := json.Marshal(map[string]interface{}{
		"connectID": connectID,
		"type":      "connect",
	})
	handle(datagram.Datagram{
		TagName:     "tcp-server",
		Addr:        addr,
		ProjectName: projectName,
		ContextID:   contextID,
		ExtraInfo:   string(extraInfo),
		Time:        time.Now().UnixNano(),
	})

	parser := ParserFactory(func(b []byte) {
		if fistPkg {
			if b[0] == '$' {
				for i := 2; i < len(b); i++ {
					if b[i] == '#' {
						if projectName == "" {
							projectName = string(b[1:i])
						} else if contextID == "" {
							contextID = string(b[1+len(projectName)+1 : i])
						} else {
							break
						}
					}
				}
				if len(projectName) > 0 {
					extraInfo, _ := json.Marshal(map[string]interface{}{
						"connectID": connectID,
						"type":      "connect",
					})
					handle(datagram.Datagram{
						TagName:     "tcp-server",
						Addr:        addr,
						ProjectName: projectName,
						ContextID:   contextID,
						ExtraInfo:   string(extraInfo),
						Time:        time.Now().UnixNano(),
					})
					return
				}
			}
			fistPkg = false
		}

		extraInfo, _ := json.Marshal(map[string]interface{}{
			"connectID": connectID,
			"type":      "datagram",
		})
		handle(datagram.Datagram{
			TagName:     "tcp-server",
			Addr:        addr,
			ProjectName: projectName,
			ContextID:   contextID,
			Content:     string(b),
			ExtraInfo:   string(extraInfo),
			Time:        time.Now().UnixNano(),
		})
	})

	for {
		var buf [128]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			reason := ""
			if err == io.EOF {
				reason = "io.EOF"
			} else {
				reason = err.Error()
			}
			extraInfo, _ := json.Marshal(map[string]interface{}{
				"connectID": connectID,
				"type":      "disconnect",
				"reason":    reason,
			})
			handle(datagram.Datagram{
				TagName:     "tcp-server",
				Addr:        addr,
				ProjectName: projectName,
				ContextID:   contextID,
				ExtraInfo:   string(extraInfo),
				Time:        time.Now().UnixNano(),
			})
			break
		}
		if n > 0 {
			parser(buf[:n])
		}
	}
}

// Listen *
func (server *tcpServer) Listen(handle func(datagram.Datagram)) {
	for {
		// 建立socket连接
		conn, err := server.listener.Accept()
		if err != nil {
			panic(fmt.Errorf("Listen.Accept failed, %v", err))
		}
		// 业务处理逻辑
		go handleConn(conn, handle)
	}
}

func (server *tcpServer) Close() {
	// TODO
}
