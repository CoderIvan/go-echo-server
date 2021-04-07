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

func format(data datagram.Datagram) []interface{} {
	params := []interface{}{
		time.Unix(0, data.Time).Format("20060102 15:04:05.000 +0800"),
		">>",
		data.TagName,
		">>",
		data.Addr,
	}

	if len(data.ProjectName) > 0 {
		params = append(params, ">>", data.ProjectName)
	}

	params = append(params, ">>", string(data.Content))

	return params
}

// Handle *
func (l *logger) Handle(data datagram.Datagram) {
	fmt.Println(format(data))
}
