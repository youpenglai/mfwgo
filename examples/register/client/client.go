package main

import (
	"flag"
	"log"
	"time"

	. "github.com/youpenglai/mfwgo/registry"
)

const (
	defaultName = "defaultService"
)

func main() {
	var serviceName = flag.String("s", defaultName, "input service name")
	flag.Parse()

	for i := 0; i < 4; i++ {
		resultInfo, err := DiscoverService(*serviceName)
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Printf("get service success, sericeInfo: %+v\n", resultInfo)
		time.Sleep(time.Second)
	}
}
