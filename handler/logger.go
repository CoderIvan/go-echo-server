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

	if len(data.HexContent) > 0 {
		params = append(params, ">>", data.HexContent)
	}

	if len(data.Content) > 0 {
		params = append(params, ">>", data.Content)
	}

	if len(data.ContextID) > 0 {
		params = append(params, ">>", data.ContextID)
	}

	if len(data.ExtraInfo) > 0 {
		params = append(params, ">>", data.ExtraInfo)
	}

	return params
}

// Handle *
func (l *logger) Handle(data datagram.Datagram) {
	fmt.Println(format(data)...)
}
