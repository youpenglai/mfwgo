package registry

import (
	"errors"
)

type ServiceRegistration struct {
	ServiceName string
	Port        int64
}

type ServiceRegistrType struct {
	CheckHealth CheckHealth
}

type CheckHealth struct {
	Type string
}

type ServiceInfo struct {
	//ID      string   `json:"ID"`
	//Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Port    int64    `json:"Port"`
	Tags    []string `json:"Tags"`
}

func RegisterService(serviceInfo ServiceRegistration, serviceType ServiceRegistrType) error {
	return NewConsulService().Register(serviceInfo.ServiceName, serviceInfo.Port, serviceType.CheckHealth.Type)
}

func DiscoverService(serviceName string) (*ServiceInfo, error) {
	service := gConsulCache.Get(serviceName)
	if service == nil {
		err := NewConsulService().GetServicesToCache(serviceName)
		if err != nil {
			return nil, err
		}

		service = gConsulCache.Get(serviceName)
		if service == nil {
			return service, errors.New("Not Found Service: " + serviceName)
		}
	}
	return service, nil
}
