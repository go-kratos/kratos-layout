package service

import (
	"context"
	"log"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// GreeterService is a greeter service.
type GreeterService struct {
	pb.UnimplementedGreeterServer
}

// NewGreeterService new a greeter service.
func NewGreeterService() *GreeterService {
	return &GreeterService{}
}

// SayHello implements helloworld.GreeterServer
func (s *GreeterService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
