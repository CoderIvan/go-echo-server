package server

import (
	"encoding/json"
	"fmt"
	"go-echo-server/datagram"
	"net"
	"strings"
	"time"
)

// UDPServer *
type udpServer struct {
	port int
}

func NewUDPServer(port int) *udpServer {
	return &udpServer{
		port: port,
	}
}

func processLed(buf []byte) (bool, string, []byte) {
	var isSuccess bool
	var projectName string
	var content []byte
	valueDotKey := strings.Split(string(buf), "=")

	if len(valueDotKey) > 1 {
		jsonMap := make(map[string]interface{})
		key := valueDotKey[0]
		for i := 1; i < len(valueDotKey)-1; i++ {
			entry := valueDotKey[i]
			lastIndex := strings.LastIndex(entry, ",")

			var v interface{} = entry[:lastIndex]
			var mapResult map[string]interface{}
			if err := json.Unmarshal([]byte(v.(string)), &mapResult); err == nil {
				v = mapResult
			}

			if key == "project" {
				projectName = v.(string)
			} else {
				jsonMap[key] = v
			}
			key = entry[lastIndex+1:]
		}
		jsonMap[key] = valueDotKey[len(valueDotKey)-1]

		if b, err := json.Marshal(jsonMap); err == nil {
			content = b
			isSuccess = true
		}
	}
	return isSuccess, projectName, content
}

func toHexString(originList []byte) string {
	var contentList []string = make([]string, len(originList))
	for i := 0; i < len(originList); i++ {
		contentList[i] = fmt.Sprintf("%02X", originList[i])
	}
	return strings.Join(contentList, " ")
}

func process(buf []byte, addr string) datagram.Datagram {
	projectName := ""
	content := buf

	if buf[0] == '$' {
		for i := 2; i < 32; i++ {
			if buf[i] == '#' {
				projectName = string(buf[1:i])
				content = buf[i+1:]
				break
			}
		}
	} else { // 对显示屏终端进行特殊处理
		if ok, newProjectName, newContent := processLed(buf); ok {
			projectName = newProjectName
			content = newContent
		}
	}

	return datagram.Datagram{
		TagName:     "udp-server",
		Addr:        addr,
		ProjectName: projectName,
		HexContent:  toHexString(content),
		Content:     string(content),
		Time:        time.Now().UnixNano(),
	}
}

// Listen *
func (server *udpServer) Listen(handle func(datagram.Datagram)) error {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: server.port,
	})
	if err != nil {
		return err
	}

	for {
		// DOTO 有待优化，这里会不断的分配内存与回收内存，可以考虑使用缓冲池
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		if n > 0 {
			go handle(process(buf[0:n], addr.String()))
		}
	}
}
