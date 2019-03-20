# mfwgo
Penglai MicroFramework golang version.

## Quick Start

#### Download and install

    go get -u github.com/youpenglai/mfwgo

Documentation
-------------
See examples find examples in the [examples directory](examples/).

registry:
-------------
```go
    registry.RegisterService(serviceInfo, serviceOptions)
    registry.DiscoverService(serviceName)
```

service:
-------------
```go
    grpcServer := server.NewGRPCServer(server.GRPCServerOption{IpAddr: ipAddr, Port: port})
    pb.RegisterGreeterServer(grpcServer.GetServer(), &server{})
    grpcServer.ListenAndServe()
```

client:
-------------
```go
    client, err := NewGRPCConn(serviceName)
    c := pb.NewGreeterClient(client.GetConn())
    c.SayHello(ctx, &pb.HelloRequest{Name: defaultName})
```