package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	grpctransport "github.com/go-kratos/kratos/v2/transport/grpc"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/service-layout/api/helloworld"
	"github.com/go-kratos/service-layout/internal/service"

	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Version is the version of the compiled software.
	Version string
	// Branch is current branch name the code is built off.
	Branch string
	// Revision is the short commit hash of source tree.
	Revision string
	// BuildDate is the date when the binary was built.
	BuildDate string
)

func main() {
	log.Printf("service version: %s\n", Version)

	// transport server
	httpSrv := httptransport.NewServer(":8000")
	grpcSrv := grpctransport.NewServer(":9000")

	// register service
	gs := service.NewGreeterService()
	helloworld.RegisterGreeterServer(grpcSrv, gs)
	helloworld.RegisterGreeterHTTPServer(httpSrv, gs)

	// application lifecycle
	app := kratos.New()
	app.Append(kratos.Hook{OnStart: httpSrv.Start, OnStop: httpSrv.Stop})
	app.Append(kratos.Hook{OnStart: grpcSrv.Start, OnStop: grpcSrv.Stop})

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Printf("startup failed: %v\n", err)
	}
}
