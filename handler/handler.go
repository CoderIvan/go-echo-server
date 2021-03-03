package handler

import "go-echo-server/datagram"

// Handler *
type Handler interface {
	Handle(data datagram.Datagram)
}
