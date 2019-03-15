package main

import (
	"fmt"
	"github.com/youpenglai/mfwgo/consul"

	"log"
)

func main() {
	err := consul.RegisterService("facelabs", 10150)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("service register success\n")

	result_info, err := consul.DiscoverService("facelabs")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("get service success, serice: %v\n", result_info)
}
