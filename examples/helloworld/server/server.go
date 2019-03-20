package main

import (
	"context"
	"log"

	pb "github.com/youpenglai/mfwgo/examples/helloworld/proto"
	. "github.com/youpenglai/mfwgo/registry"
	"github.com/youpenglai/mfwgo/server"
)

const (
	serviceName = "hello world"
	ipAddr      = "127.0.0.1"
	port        = 20003
)

type serverTest struct{}

func (s *serverTest) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	//time.Sleep(50 * time.Second)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// register service
	err := RegisterService(ServiceRegistration{ServiceName: serviceName, Port: port},
		ServiceRegisterType{CheckHealth: CheckHealth{Type: "grpc"}})
	if err != nil {
		log.Fatalf("register service error: %v", err)
	}

	// start server
	grpcServer := server.NewGRPCServer(server.GRPCServerOption{IpAddr: ipAddr, Port: port})
	pb.RegisterGreeterServer(grpcServer.GetServer(), &serverTest{})
	grpcServer.ListenAndServe()
}
