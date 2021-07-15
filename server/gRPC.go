package server

import (
	"context"
	"go-echo-server/datagram"
	pb "go-echo-server/server/message"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type serv struct {
	port int
}

func NewGRPCServer(port int) *serv {
	return &serv{
		port: port,
	}
}

func (this *serv) Listen(handle func(datagram.Datagram)) error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(this.port))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &gServer{
		handle: handle,
	})
	return s.Serve(lis)
}

type gServer struct {
	pb.UnimplementedGreeterServer

	handle func(datagram.Datagram)
}

func (s *gServer) SendText(ctx context.Context, in *pb.TextRequest) (*pb.TextResponse, error) {
	s.handle(datagram.Datagram{
		TagName:     "gRPC",
		ProjectName: in.GetProjectName(),
		Content:     in.GetMessage(),
		Time:        time.Now().UnixNano(),
	})
	return &pb.TextResponse{Message: "success"}, nil
}
