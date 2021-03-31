package handler

import (
	"fmt"
	"go-echo-server/datagram"
	"time"
)

// Logger *
type logger struct {
}

func NewLogger() *logger {
	return &logger{}
}

// Handle *
func (l *logger) Handle(data datagram.Datagram) {
	if len(data.ProjectName) > 0 {
		fmt.Println(time.Unix(0, data.Time).Format("20060102 15:04:05.000 +0800"), ">>", data.TagName, ">>", data.Addr, ">>", data.ProjectName, ">>", string(data.Content))
	} else {
		fmt.Println(time.Unix(0, data.Time).Format("20060102 15:04:05.000 +0800"), ">>", data.TagName, ">>", data.Addr, ">>", string(data.Content))
	}
}
