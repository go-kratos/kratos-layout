package main

import (
	"flag"
	"os"

	pb "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/go-kratos/kratos-layout/internal/data"
	"github.com/go-kratos/kratos-layout/internal/server"
	"github.com/go-kratos/kratos-layout/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/log/stdlog"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"gopkg.in/yaml.v2"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newApp(logger log.Logger, hs *http.Server, gs *grpc.Server, greeter *service.GreeterService) *kratos.App {
	pb.RegisterGreeterServer(gs, greeter)
	pb.RegisterGreeterHTTPServer(hs, greeter)
	return kratos.New(
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	flag.Parse()
	logger := stdlog.NewLogger(stdlog.Writer(os.Stdout))
	defer logger.Close()

	conf := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(func(kv *config.KeyValue, v map[string]interface{}) error {
			return yaml.Unmarshal(kv.Value, v)
		}),
		config.WithLogger(logger),
	)
	if err := conf.Load(); err != nil {
		panic(err)
	}

	var (
		hc server.HTTPConfig
		gc server.GRPCConfig
		dc data.Config
	)
	if err := conf.Value("http.server").Scan(&hc); err != nil {
		panic(err)
	}
	if err := conf.Value("grpc.server").Scan(&gc); err != nil {
		panic(err)
	}
	if err := conf.Value("data").Scan(&dc); err != nil {
		panic(err)
	}

	app, err := InitApp(&hc, &gc, &dc, logger)
	if err != nil {
		panic(err)
	}

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
