package client

import (
	"fmt"

	"github.com/youpenglai/mfwgo/registry"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	cc          *grpc.ClientConn
	serviceInfo *registry.ServiceInfo
}

func NewGRPCConn(serviceName string) (*GRPCClient, error) {
	s, err := registry.DiscoverService(serviceName)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", s.Address, s.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		cc:          conn,
		serviceInfo: s,
	}, nil
}

func (this *GRPCClient) GetConn() *grpc.ClientConn {
	return this.cc
}
