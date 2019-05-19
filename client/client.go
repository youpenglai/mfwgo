package client

import (
	"fmt"

	"google.golang.org/grpc"
)


type GRPCClient struct {
	cc          *grpc.ClientConn
}

func getServiceUrl(serviceName string) string {
	return fmt.Sprintf("%s:///%s", mfwScheme, serviceName)
}

func NewGRPCConn(serviceName string) (*GRPCClient, error) {
	conn, err := grpc.Dial(
		getServiceUrl(serviceName),
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"), // 使用轮询调度
	)
	if err != nil {
		return nil, err
	}

	return &GRPCClient{
		cc:          conn,
	}, nil
}

func (c *GRPCClient) GetConn() *grpc.ClientConn {
	return c.cc
}
