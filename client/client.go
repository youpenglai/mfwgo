package client

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/youpenglai/mfwgo/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var (
	grpcStateCheckTime = 10 * time.Second
)

type GRPCClient struct {
	cc *grpc.ClientConn
	ch chan bool
}

func getServiceUrl(serviceName string) string {
	return fmt.Sprintf("%s:///%s", mfwScheme, serviceName)
}

func (this *GRPCClient) Alter() <-chan bool {
	return this.ch
}

func NewGRPCConn(serviceName string) (*GRPCClient, error) {
	cc, err := newConn(serviceName)
	if err != nil {
		return nil, err
	}

	client := GRPCClient{
		cc: cc,
		ch: make(chan bool),
	}

	go func() {
		for range time.Tick(grpcStateCheckTime) {
			switch client.cc.GetState() {
			case connectivity.Idle, connectivity.TransientFailure:
				cc, err := reConnect(serviceName)
				if err != nil {
					log.Printf("err: %v", err.Error())
					continue
				}
				client.cc = cc
				client.ch <- true
			}
		}
	}()

	return &client, nil
}

func reConnect(serviceName string) (*grpc.ClientConn, error) {
	services, err := registry.NewConsulService().GetServices(serviceName)
	if err != nil {
		return nil, err
	}
	if len(services) == 0 {
		return nil, errors.New("Not found service: " + serviceName)
	}

	return newConn(serviceName)
}

func newConn(serviceName string) (*grpc.ClientConn, error) {
	return grpc.Dial(
		getServiceUrl(serviceName),
		grpc.WithInsecure(),
		grpc.WithBalancerName("round_robin"), // 使用轮询调度
	)
}

func (c *GRPCClient) GetConn() *grpc.ClientConn {
	return c.cc
}
