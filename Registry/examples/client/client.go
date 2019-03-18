package main

import (
	"fmt"
	"log"

	"github.com/youpenglai/mfwgo/Registry/consul"
)

func main() {

	httpRegister("facelabs", 20001)
	grpcRegister("facelabs", 20002)

	result_info, err := consul.DiscoverService("facelabs")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("get service success, serice: %v\n", result_info)

	result_info, err = consul.DiscoverService("facelabs")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("get service success, serice: %v\n", result_info)
}

func httpRegister(serviceName string, servicePort int64) {
	err := consul.RegisterService(&consul.ServiceRegistration{serviceName, servicePort},
		&consul.ServiceRegistrType{CheckHealth:consul.CheckHealth{"http"}})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("service [%s:%d] register success\n\n", serviceName, servicePort)
}

func grpcRegister(serviceName string, servicePort int64) {
	err := consul.RegisterService(&consul.ServiceRegistration{serviceName, servicePort},
		&consul.ServiceRegistrType{CheckHealth:consul.CheckHealth{"grpc"}})
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("service [%s:%d] register success\n\n", serviceName, servicePort)
}
