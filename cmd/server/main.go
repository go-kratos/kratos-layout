package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-kratos/kratos/v2"
	grpctransport "github.com/go-kratos/kratos/v2/transport/grpc"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/service-layout/internal/service"

	"google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	log.Printf("version: %s\n", Version)

	// transport server
	httpSrv := httptransport.NewServer(httptransport.WithAddress(":8000"))
	grpcSrv := grpctransport.NewServer(grpctransport.WithAddress(":9000"))

	// register service
	gs := service.NewGreeterService()
	helloworld.RegisterGreeterServer(grpcSrv, gs)
	httpSrv.Handler = newHTTPHandler(gs)

	// application lifecycle
	app := kratos.New()
	app.Append(kratos.Hook{OnStart: httpSrv.Start, OnStop: httpSrv.Stop})
	app.Append(kratos.Hook{OnStart: grpcSrv.Start, OnStop: grpcSrv.Stop})

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Printf("app failed: %v\n", err)
	}
}

func newHTTPHandler(gs *service.GreeterService) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/helloworld", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "helloworld")
	})
	return m
}
