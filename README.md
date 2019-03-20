# mfwgo
Penglai MicroFramework golang version.

Download and install
-------------

    go get -u github.com/youpenglai/mfwgo

Documentation
-------------
See examples find examples in the [examples directory](examples/).

registry:
-------------
```go
    // register service
    registry.RegisterService(serviceInfo, serviceOptions)
    // get service
    registry.DiscoverService(serviceName)
```

service:
-------------
```go
    grpcServer := server.NewGRPCServer(server.GRPCServerOption{IpAddr: ipAddr, Port: port})
    pb.RegisterGreeterServer(grpcServer.GetServer(), &serverTest{})
    grpcServer.ListenAndServe()
```

client:
-------------
```go
    cc, err := client.NewGRPCConn(serviceName)
    c := pb.NewGreeterClient(cc.GetConn())
    c.SayHello(ctx, &pb.HelloRequest{Name: defaultName})
```