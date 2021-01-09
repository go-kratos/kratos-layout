package main

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/log/stdlog"
	servergrpc "github.com/go-kratos/kratos/v2/server/grpc"
	serverhttp "github.com/go-kratos/kratos/v2/server/http"
	grpctransport "github.com/go-kratos/kratos/v2/transport/grpc"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/service-layout/api/helloworld"
	"github.com/go-kratos/service-layout/internal/service"
	"google.golang.org/grpc"

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
	logger, err := stdlog.NewLogger(stdlog.Writer(os.Stdout))
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	log := log.NewHelper("main", logger)
	log.Infof("version: %s", Version)

	// transport
	httpTransport := httptransport.NewServer()
	grpcTransport := grpctransport.NewServer()

	// server
	httpServer := serverhttp.NewServer("tcp", ":8000", serverhttp.ServerHandler(httpTransport))
	grpcServer := servergrpc.NewServer("tcp", ":9000", grpc.UnaryInterceptor(grpcTransport.Interceptor()))

	// register service
	gs := service.NewGreeterService()
	helloworld.RegisterGreeterServer(grpcServer, gs)
	helloworld.RegisterGreeterHTTPServer(httpTransport, gs)

	// application lifecycle
	app := kratos.New()
	app.Append(httpServer)
	app.Append(grpcServer)

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		log.Errorf("start failed: %v\n", err)
	}
}
