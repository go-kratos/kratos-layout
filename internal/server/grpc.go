package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/util/time"
)

// GRPCConfig is server config.
type GRPCConfig struct {
	Network string        `json:"network"`
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
}

// NewGRPCServer new a gRPC server.
func NewGRPCServer(c *GRPCConfig) *grpc.Server {
	var opts []grpc.ServerOption
	if c.Network != "" {
		opts = append(opts, grpc.Network(c.Network))
	}
	if c.Address != "" {
		opts = append(opts, grpc.Address(c.Address))
	}
	if c.Timeout != 0 {
		opts = append(opts, grpc.Timeout(c.Timeout.Duration()))
	}
	return grpc.NewServer(opts...)
}
