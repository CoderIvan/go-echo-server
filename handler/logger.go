package handler

import (
	"fmt"
	"go-echo-server/datagram"
	"time"
)

// Logger *
type Logger struct {
}

// Handle *
func (l *Logger) Handle(data datagram.Datagram) {
	if len(data.ProjectName) > 0 {
		fmt.Println(time.Now().Format("20060102 15:04:05.000 +0800"), ">>", data.TagName, ">>", data.Addr, ">>", data.ProjectName, ">>", string(data.Content))
	} else {
		fmt.Println(time.Now().Format("20060102 15:04:05.000 +0800"), ">>", data.TagName, ">>", data.Addr, ">>", string(data.Content))
	}
}
