package common

import "google.golang.org/grpc"

type GRPCConfig struct {
	Port int `mapstructure:"port" yaml:"port"`
}

func (c *GRPCConfig) GetPort() int {
	if c.Port == 0 {
		return 10000
	}
	return c.Port
}

type GRPCClientConfig struct {
	Insecure bool   `mapstructure:"insecure" yaml:"insecure"`
	Addr     string `mapstructure:"addr" yaml:"addr"`
}

func (c *GRPCClientConfig) Connect() (*grpc.ClientConn, error) {
	var opts []grpc.DialOption
	if c.Insecure {
		opts = append(opts, grpc.WithInsecure())
	}

	return grpc.Dial(c.Addr, opts...)
}
