package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type TCPListenPort int

func NewServer() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}

type ClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func NewClient(c ClientConfig) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.Insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return grpc.Dial(c.Addr, opts...)
}
