package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPConfig is server config.
type HTTPConfig struct {
	Network string `json:"network"`
	Address string `json:"address"`
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *HTTPConfig, logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Logger(logger),
	}
	if c.Network != "" {
		opts = append(opts, http.Network(c.Network))
	}
	if c.Address != "" {
		opts = append(opts, http.Address(c.Address))
	}
	return http.NewServer(opts...)
}
