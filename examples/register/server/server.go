package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	pb "github.com/youpenglai/mfwgo/examples/register/proto"
	. "github.com/youpenglai/mfwgo/registry"

	"google.golang.org/grpc"
)

func main() {
	ch := make(chan error, 1)
	go httpDo("my-service", "127.0.0.1", 20001, ch)
	go grpcDo("my-service", "127.0.0.1", 20002, ch)
	log.Fatalf("err: %v", (<-ch).Error())
}

func httpDo(name string, address string, port int64, out chan<- error) {
	err := RegisterService(ServiceRegistration{ServiceName: name, Port: port},
		ServiceRegistrType{CheckHealth: CheckHealth{Type: "http"}})
	if err != nil {
		out <- err
		return
	}
	log.Printf("service [%s:%d] http register success\n\n", name, port)

	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	out <- http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), nil)
}

type server struct{}

func (s *server) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func grpcDo(name string, address string, port int64, out chan<- error) {
	err := RegisterService(ServiceRegistration{ServiceName: name, Port: port},
		ServiceRegistrType{CheckHealth: CheckHealth{Type: "grpc"}})
	if err != nil {
		out <- err
		return
	}
	log.Printf("service [%s:%d] grpc register success\n\n", name, port)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, port))
	if err != nil {
		out <- err
		return
	}

	s := grpc.NewServer()
	pb.RegisterHealthServer(s, &server{})
	out <- s.Serve(lis)
}
