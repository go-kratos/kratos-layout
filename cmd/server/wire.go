// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/data"
	"github.com/go-kratos/kratos-layout/internal/server"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// InitApp init kratos application.
func InitApp(*server.HTTPConfig, *server.GRPCConfig, *data.Config, log.Logger) (*kratos.App, error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, biz.ProviderSet, service.ProviderSet, newApp))
}
