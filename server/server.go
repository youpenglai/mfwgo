package server

import (
	"context"
	"fmt"
	"net"

	pb "github.com/youpenglai/mfwgo/server/proto"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server  *grpc.Server
	Options GRPCServerOption
}

type GRPCServerOption struct {
	IpAddr string
	Port   int64
}

type server struct{}

func (s *server) Check(ctx context.Context, in *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{Status: pb.HealthCheckResponse_SERVING}, nil
}

func NewGRPCServer(options GRPCServerOption) *GRPCServer {
	return &GRPCServer{
		server:  grpc.NewServer(),
		Options: options,
	}
}

func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

func (s *GRPCServer) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Options.IpAddr, s.Options.Port))
	if err != nil {
		return err
	}

	pb.RegisterHealthServer(s.server, &server{})
	return s.server.Serve(lis)
}
