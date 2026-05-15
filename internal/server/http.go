package server

import (
	"log/slog"

	hellov1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	routev1 "github.com/go-kratos/kratos-layout/api/routeguide/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v3/middleware/recovery"
	"github.com/go-kratos/kratos/v3/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, routeGuide *service.RouteGuideService, logger *slog.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	hellov1.RegisterGreeterHTTPServer(srv, greeter)
	routev1.RegisterRouteGuideHTTPServer(srv, routeGuide)
	return srv
}
