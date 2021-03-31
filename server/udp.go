package server

import (
	"encoding/json"
	"go-echo-server/datagram"
	"net"
	"strings"
	"time"
)

// UDPServer *
type udpServer struct {
	conn *net.UDPConn
}

func NewUDPServer(port int) *udpServer {
	conn, _ := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
	})
	return &udpServer{
		conn: conn,
	}
}

func processLed(buf []byte) (bool, string, []byte) {
	var projectName string
	var content []byte
	valueDotKey := strings.Split(string(buf), "=")

	if len(valueDotKey) > 1 {
		keys := make([]string, 0, len(valueDotKey))
		values := make([]interface{}, 0, len(valueDotKey))

		keys = append(keys, valueDotKey[0])
		for i := 1; i < len(valueDotKey)-1; i++ {
			entry := valueDotKey[i]
			lastIndex := strings.LastIndex(entry, ",")
			var v interface{} = entry[:lastIndex]
			var mapResult map[string]interface{}
			if err := json.Unmarshal([]byte(v.(string)), &mapResult); err == nil {
				v = mapResult
			}
			k := entry[lastIndex+1:]
			keys = append(keys, k)
			values = append(values, v)
		}
		values = append(values, valueDotKey[len(valueDotKey)-1])

		jsonMap := make(map[string]interface{})
		for i := 0; i < len(keys); i++ {
			if keys[i] == "project" {
				projectName = values[i].(string)
			} else {
				jsonMap[keys[i]] = values[i]
			}
		}
		if b, err := json.Marshal(jsonMap); err == nil {
			content = b
		}
		return true, projectName, content
	}
	return false, "", nil
}

func process(buf []byte, addr string) datagram.Datagram {
	projectName := ""
	content := buf

	if buf[0] == '$' {
		for i := 2; i < 32; i++ {
			if buf[i] == '#' {
				projectName = string(buf[1:i])
				content = buf[i+1:]
			}
		}
	} else { // 对显示屏终端进行特殊处理
		ok, newProjectName, newContent := processLed(buf)

		if ok {
			projectName = newProjectName
			content = newContent
		}
	}

	return datagram.Datagram{
		TagName:     "udp-server",
		Addr:        addr,
		ProjectName: projectName,
		Content:     string(content),
		Time:        time.Now().UnixNano(),
	}
}

// Listen *
func (server *udpServer) Listen(handle func(datagram.Datagram)) {
	for {
		// DOTO 有待优化，这里会不断的分配内存与回收内存，可以考虑使用缓冲池
		buf := make([]byte, 1024)
		n, addr, _ := server.conn.ReadFromUDP(buf)

		if n > 0 {
			handle(process(buf[0:n], addr.String()))
		}
	}
}

func (server *udpServer) Close() {
	server.conn.Close()
}
