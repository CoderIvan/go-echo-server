package server

import "go-echo-server/datagram"

// Server *
type Server interface {
	Listen(ch chan datagram.Datagram)
}
