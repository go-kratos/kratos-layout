package helloworld

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-kratos/kratos/v2/encoding"
	httptransport "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/mux"
)

// RegisterGreeterHTTPServer register http server.
func RegisterGreeterHTTPServer(srv *httptransport.Server, si GreeterServer) {
	s := &greeterHTTPServer{opts: srv.Options(), service: si}
	r := mux.NewRouter()
	r.HandleFunc("/helloworld", s.SayHello).Methods("POST")
	srv.AddHandler(r)
}

type greeterHTTPServer struct {
	opts    httptransport.ServerOptions
	service GreeterServer
}

func (s *greeterHTTPServer) SayHello(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("content-type")
	name, _ := httptransport.ContentSubtype(contentType)
	codec := encoding.GetCodec(name)
	fmt.Println(name, contentType)
	if codec == nil {
		s.opts.ErrorHandler(req.Context(), httptransport.ErrUnknownCodec(contentType), codec, res)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.opts.ErrorHandler(req.Context(), httptransport.ErrDataLoss(err.Error()), codec, res)
		return
	}
	defer req.Body.Close()
	in := new(HelloRequest)
	if err = codec.Unmarshal(data, in); err != nil {
		s.opts.ErrorHandler(req.Context(), httptransport.ErrCodecUnmarshal(err.Error()), codec, res)
		return
	}
	out, err := s.service.SayHello(req.Context(), in)
	if err != nil {
		s.opts.ErrorHandler(req.Context(), err, codec, res)
		return
	}
	body, err := codec.Marshal(out)
	if err != nil {
		s.opts.ErrorHandler(req.Context(), httptransport.ErrCodecMarshal(err.Error()), codec, res)
		return
	}
	res.Write(body)
}
