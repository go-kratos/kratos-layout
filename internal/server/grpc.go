package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// GRPCConfig is server config.
type GRPCConfig struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *GRPCConfig, logger log.Logger) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Logger(logger),
	}
	if c.Network != "" {
		opts = append(opts, grpc.Network(c.Network))
	}
	if c.Address != "" {
		opts = append(opts, grpc.Address(c.Address))
	}
	return grpc.NewServer(opts...)
}
