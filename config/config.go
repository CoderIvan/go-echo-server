package config

import (
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

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
		TCP struct {
			PORT int `yaml:"port"`
		} `yaml:"tcp"`
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

// Get *
func Get() (config, error) {
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

	if os.Getenv("SERVER_TCP_PORT") != "" {
		c.SERVER.TCP.PORT, _ = strconv.Atoi(os.Getenv("SERVER_TCP_PORT"))
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
