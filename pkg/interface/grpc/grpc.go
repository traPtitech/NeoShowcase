package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type TCPListenPort int

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}
