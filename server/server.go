package server

import "go-echo-server/handler"

// Server *
type Server interface {
	Listen(h handler.Handler)
}
