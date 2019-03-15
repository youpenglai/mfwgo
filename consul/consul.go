package consul

import "errors"

func RegisterService(serviceName string, port int64) error {
	return NewConsulService().Register(serviceName, port)
}

func DiscoverService(serviceName string) (*ServiceInfo, error) {
	service := g_consulCache.Get(serviceName)
	if service == nil {
		return service, errors.New("Not Found Service: " + serviceName)
	}
	return service, nil
}
