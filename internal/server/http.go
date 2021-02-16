package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/kratos/v2/util/time"
)

// HTTPConfig is server config.
type HTTPConfig struct {
	Network string        `json:"network"`
	Address string        `json:"address"`
	Timeout time.Duration `json:"timeout"`
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(c *HTTPConfig) *http.Server {
	var opts []http.ServerOption
	if c.Network != "" {
		opts = append(opts, http.Network(c.Network))
	}
	if c.Address != "" {
		opts = append(opts, http.Address(c.Address))
	}
	if c.Timeout != 0 {
		opts = append(opts, http.Timeout(c.Timeout.Duration()))
	}
	return http.NewServer(opts...)
}
