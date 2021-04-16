package main

import (
	"fmt"
	"go-echo-server/datagram"
	"go-echo-server/handler"
	"go-echo-server/server"
	"io/ioutil"
	"os"
	"strconv"
	"sync"

	"gopkg.in/yaml.v3"
)

type connect struct {
	servers  []server.Server
	handlers []handler.Handler
}

func (ct connect) run() {
	var wg sync.WaitGroup
	for _, server := range ct.servers {
		wg.Add(1)
		go server.Listen(func(data datagram.Datagram) {
			var wg2 sync.WaitGroup
			for _, hand := range ct.handlers {
				wg2.Add(1)
				go func(hd handler.Handler) {
					hd.Handle(data)
					wg2.Done()
				}(hand)
			}
			wg2.Wait()
		})
	}
	wg.Wait()
}

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
// Config *
type config struct {
	SERVER struct {
		UDP struct {
			PORT int `yaml:"port"`
		} `yaml:"udp"`
		HTTP struct {
			PORT int `yaml:"port"`
		} `yaml:"http"`
		MQTT struct {
			PORT int `yaml:"port"`
		} `yaml:"mqtt"`
	} `yaml:"server"`
	HANDLER struct {
		SLS struct {
			ACCESSKEYID     string `yaml:"accessKeyID"`
			ACCESSKEYSECRET string `yaml:"accessKeySecret"`
			ENDPOINT        string `yaml:"endpoint"`
			PROJECTNAME     string `yaml:"projectName"`
			LOGSTORENAME    string `yaml:"logStoreName"`
		} `yaml:"sls"`
	} `yaml:"handler"`
}

// GetConfig *
func getConfig() (config, error) {
	c := config{}

	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal([]byte(content), &c)
	if err != nil {
		return c, err
	}

	if os.Getenv("SERVER_UDP_PORT") != "" {
		c.SERVER.UDP.PORT, _ = strconv.Atoi(os.Getenv("SERVER_UDP_PORT"))
	}

	if os.Getenv("SERVER_HTTP_PORT") != "" {
		c.SERVER.HTTP.PORT, _ = strconv.Atoi(os.Getenv("SERVER_HTTP_PORT"))
	}

	if os.Getenv("SERVER_MQTT_PORT") != "" {
		c.SERVER.MQTT.PORT, _ = strconv.Atoi(os.Getenv("SERVER_MQTT_PORT"))
	}

	if os.Getenv("HANDLER_SLS_ACCESSKEYID") != "" {
		c.HANDLER.SLS.ACCESSKEYID = os.Getenv("HANDLER_SLS_ACCESSKEYID")
	}

	if os.Getenv("HANDLER_SLS_ACCESSKEYSECRET") != "" {
		c.HANDLER.SLS.ACCESSKEYSECRET = os.Getenv("HANDLER_SLS_ACCESSKEYSECRET")
	}

	if os.Getenv("HANDLER_SLS_ENDPOINT") != "" {
		c.HANDLER.SLS.ENDPOINT = os.Getenv("HANDLER_SLS_ENDPOINT")
	}

	if os.Getenv("HANDLER_SLS_PROJECTNAME") != "" {
		c.HANDLER.SLS.PROJECTNAME = os.Getenv("HANDLER_SLS_PROJECTNAME")
	}

	if os.Getenv("HANDLER_SLS_LOGSTORENAME") != "" {
		c.HANDLER.SLS.LOGSTORENAME = os.Getenv("HANDLER_SLS_LOGSTORENAME")
	}

	return c, nil
}

func main() {
	config, err := getConfig()

	if err != nil {
		fmt.Println("读取配置错误", err)
		return
	}

	ct := connect{
		[]server.Server{
			server.NewUDPServer(config.SERVER.UDP.PORT),
			server.NewHTTPServer(config.SERVER.HTTP.PORT),
			server.NewMqttServer(config.SERVER.MQTT.PORT),
		},
		[]handler.Handler{
			handler.NewLogger(),
			handler.NewSLS(
				config.HANDLER.SLS.ACCESSKEYID,
				config.HANDLER.SLS.ACCESSKEYSECRET,
				config.HANDLER.SLS.ENDPOINT,
				config.HANDLER.SLS.PROJECTNAME,
				config.HANDLER.SLS.LOGSTORENAME,
			),
		},
	}

	ct.run()
}
