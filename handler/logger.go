package handler

import (
	"fmt"
	"time"
)

type Logger struct {
}

func (l *Logger) Handle(protocol string, addr string, content string, projectName string) {
	if len(projectName) > 0 {
		fmt.Println(time.Now().Format("20060102 15:04:05.000 +0800"), ">>", protocol, ">>", addr, ">>", projectName, ">>", string(content))
	} else {
		fmt.Println(time.Now().Format("20060102 15:04:05.000 +0800"), ">>", protocol, ">>", addr, ">>", string(content))
	}
}
