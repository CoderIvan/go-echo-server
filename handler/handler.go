package handler

// Handler *
type Handler interface {
	Handle(protocol string, addr string, content string, projectName string)
}
