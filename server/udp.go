package server

import (
	"encoding/json"
	"go-echo-server/datagram"
	"net"
	"strings"
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
	} else { // 对显示屏终端进行特殊处理
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
			for i := 0; i < len(keys)-1; i++ {
				if keys[i] == "project" {
					projectName = values[i].(string)
				} else {
					jsonMap[keys[i]] = values[i]
				}
			}
			if b, err := json.Marshal(jsonMap); err == nil {
				content = b
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
