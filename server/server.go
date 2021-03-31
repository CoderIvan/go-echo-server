package server

import "go-echo-server/datagram"

// Server *
type Server interface {
	Listen(handle func(datagram.Datagram))
}
