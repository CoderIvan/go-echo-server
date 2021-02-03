package handler

type HandlerI interface {
	Handle(protocol string, addr string, content string, projectName string)
}
