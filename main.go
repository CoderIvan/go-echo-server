package main

import (
	"context"
	"fmt"
	"go-echo-server/config"
	"go-echo-server/datagram"
	"go-echo-server/handler"
	"go-echo-server/server"

	"golang.org/x/sync/errgroup"
)

type connect struct {
	servers  []server.Server
	handlers []handler.Handler
}

func (ct connect) run() error {
	g, _ := errgroup.WithContext(context.Background())

	for _, server := range ct.servers {
		server := server
		g.Go(func() error {
			return server.Listen(func(data datagram.Datagram) {
				for _, hand := range ct.handlers {
					go hand.Handle(data)
				}
			})
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func main() {
	cfg, err := config.Get()
	if err != nil {
		panic(fmt.Errorf("配置读取错误, %+v", err))
	}

	ct := connect{
		[]server.Server{
			server.NewUDPServer(cfg.SERVER.UDP.PORT),
			server.NewHTTPServer(cfg.SERVER.HTTP.PORT),
			server.NewMqttServer(cfg.SERVER.MQTT.PORT),
			server.NewTCPServer(cfg.SERVER.TCP.PORT),
		},
		[]handler.Handler{
			handler.NewLogger(),
			handler.NewSLS(
				cfg.HANDLER.SLS.ACCESSKEYID,
				cfg.HANDLER.SLS.ACCESSKEYSECRET,
				cfg.HANDLER.SLS.ENDPOINT,
				cfg.HANDLER.SLS.PROJECTNAME,
				cfg.HANDLER.SLS.LOGSTORENAME,
			),
		},
	}

	if err := ct.run(); err != nil {
		panic(fmt.Errorf("服务关闭, %+v", err))
	}
}
