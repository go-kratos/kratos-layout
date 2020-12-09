package helloworld

import (
	context "context"

	"github.com/go-kratos/kratos/v2/transport/http"
)

func RegisterGreeterHTTPServer(s http.ServiceRegistrar, srv GreeterServer) {
	s.RegisterService(&_HTTP_Greeter_serviceDesc, srv)
}

func _HTTP_Greeter_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	return srv.(GreeterServer).SayHello(ctx, in)
}

var _HTTP_Greeter_serviceDesc = http.ServiceDesc{
	ServiceName: "helloworld.Greeter",
	HandlerType: (*GreeterServer)(nil),
	Methods: []http.MethodDesc{
		{
			Path:    "/helloworld",
			Method:  "POST",
			Handler: _HTTP_Greeter_SayHello_Handler,
		},
	},
	Metadata: "api/helloworld/helloworld.proto",
}
