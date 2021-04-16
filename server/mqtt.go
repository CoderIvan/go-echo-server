package server

import (
	"context"
	"encoding/json"
	"go-echo-server/datagram"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/DrmagicE/gmqtt/config"
	_ "github.com/DrmagicE/gmqtt/persistence"
	"github.com/DrmagicE/gmqtt/pkg/codes"
	"github.com/DrmagicE/gmqtt/server"
	_ "github.com/DrmagicE/gmqtt/topicalias/fifo"
)

// HTTPServer *
type mqttServer struct {
	port   int
	server *http.Server
}

func NewMqttServer(port int) *mqttServer {
	return &mqttServer{
		port: port,
	}
}

func (this *mqttServer) Listen(handle func(datagram.Datagram)) {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(this.port))
	if err != nil {
		panic(err)
	}

	var onBasicAuth server.OnBasicAuth = func(ctx context.Context, client server.Client, req *server.ConnectRequest) error {
		extraInfo, _ := json.Marshal(map[string]interface{}{
			"type": "connect",
		})

		handle(datagram.Datagram{
			TagName:     "mqtt-server",
			Addr:        client.Connection().RemoteAddr().String(),
			ProjectName: string(req.Connect.Username),
			ContextID:   string(req.Connect.ClientID),
			ExtraInfo:   string(extraInfo),
			Time:        time.Now().UnixNano(),
		})
		return nil
	}

	var onSubscribe server.OnSubscribe = func(ctx context.Context, client server.Client, req *server.SubscribeRequest) error {
		for topicName := range req.Subscriptions {
			req.Reject(topicName, &codes.Error{
				Code: codes.NotAuthorized,
				ErrorDetails: codes.ErrorDetails{
					ReasonString: []byte("publish only"),
				},
			})
		}
		return nil
	}

	var onMsgArrived server.OnMsgArrived = func(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {
		extraInfo, _ := json.Marshal(map[string]interface{}{
			"type":     "publish",
			"topic":    req.Message.Topic,
			"qos":      req.Message.QoS,
			"retained": req.Message.Retained,
		})
		handle(datagram.Datagram{
			TagName:     "mqtt-server",
			Addr:        client.Connection().RemoteAddr().String(),
			ProjectName: string(client.ClientOptions().Username),
			ContextID:   client.ClientOptions().ClientID,
			Content:     string(req.Message.Payload),
			ExtraInfo:   string(extraInfo),
			Time:        time.Now().UnixNano(),
		})
		return nil
	}

	var OnClosed server.OnClosed = func(ctx context.Context, client server.Client, err error) {
		var reason string
		if err != nil {
			reason = err.Error()
		}
		extraInfo, _ := json.Marshal(map[string]interface{}{
			"type":   "disconnect",
			"reason": reason,
		})
		data := datagram.Datagram{
			TagName:     "mqtt-server",
			Addr:        client.Connection().RemoteAddr().String(),
			ProjectName: string(client.ClientOptions().Username),
			ContextID:   client.ClientOptions().ClientID,
			ExtraInfo:   string(extraInfo),
			Time:        time.Now().UnixNano(),
		}
		handle(data)
	}

	s := server.New(
		server.WithTCPListener(ln),
		server.WithHook(server.Hooks{
			OnBasicAuth:  onBasicAuth,
			OnSubscribe:  onSubscribe,
			OnMsgArrived: onMsgArrived,
			OnClosed:     OnClosed,
		}),
		server.WithConfig(config.DefaultConfig()),
	)

	s.Run()
}

func (this *mqttServer) Close() {
}
