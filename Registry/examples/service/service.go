package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/youpenglai/mfwgo/Registry/examples/proto"
	"google.golang.org/grpc"
)

func main() {
	go httpDo()
	go grpcDo()
	time.Sleep(time.Hour * 24 * 30)
}

func httpDo() {
	http.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("ok"))
	})
	log.Fatalf("err: %v", http.ListenAndServe("127.0.0.1:20001", nil))
}

type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	log.Printf("Received: %v", in.Service)
	return &pb.HealthCheckResponse{Status: 1}, nil
}

func grpcDo() {
	lis, err := net.Listen("tcp", "127.0.0.1:20002")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterHealthServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
