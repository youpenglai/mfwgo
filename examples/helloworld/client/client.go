package main

import (
	"context"
	"log"
	"time"

	. "github.com/youpenglai/mfwgo/client"
	pb "github.com/youpenglai/mfwgo/examples/helloworld/proto"
)

const (
	serviceName = "hello world"
	defaultName = "embiid"
)

func sayHello(client pb.GreeterClient) {
	// set network call timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// call func(SayHello)
	r, err := client.SayHello(ctx, &pb.HelloRequest{Name: defaultName})
	if err != nil {
		log.Fatalf("could not greet: %v\n", err)
	}
	log.Printf("Greeting: %s\n", r.Message)
}

func main() {
	// get client connect
	client, err := NewGRPCConn(serviceName)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer client.GetConn().Close()

	c := pb.NewGreeterClient(client.GetConn())
	sayHello(c)
}
