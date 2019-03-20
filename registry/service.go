package registry

import (
	"encoding/json"
	"fmt"
)

const (
	consulIp   = "127.0.0.1"
	consulPort = 8500
)

const (
	deregisterInterval = "10m"
	tagGo              = "go"
	checkInterval      = "30s"
)

type ServiceRegisterInfo struct {
	ID    string      `json:"ID"`
	Name  string      `json:"Name"`
	Port  int64       `json:"Port"`
	Tags  []string    `json:"Tags"`
	Check interface{} `json:"Check"`
}

type HTTPCheck struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	HTTP                           string `json:"HTTP"`
	Interval                       string `json:"Interval"`
}

type GRPCCheck struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	GRPC                           string `json:"GRPC"`
	Interval                       string `json:"Interval"`
}

type ServiceHealthInfo struct {
	Node    Node              `json:"Node"`
	Service ConsulServiceInfo `json:"Service"`
	Checks  []Check           `json:"Checks"`
}

type Node struct {
	Address string `json:"Address"`
}

type Check struct {
	ServiceID   string `json:"ServiceID"`
	ServiceName string `json:"ServiceName"`
	Status      string `json:"Status"`
}

type ConsulServiceInfo struct {
	ID   string   `json:"ID"`
	Name string   `json:"Service"`
	Port int64    `json:"Port"`
	Tags []string `json:"Tags"`
}

type ConsulService struct{}

func NewConsulService() *ConsulService {
	return &ConsulService{}
}

func (c *ConsulService) Register(name string, port int64, healthType string) error {
	url := fmt.Sprintf("http://%s:%d/v1/agent/service/register", consulIp, consulPort)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	var check interface{}
	switch healthType {
	case "http":
		check = HTTPCheck{
			DeregisterCriticalServiceAfter: deregisterInterval,
			HTTP:                           fmt.Sprintf("http://%s:%d/health", consulIp, port),
			Interval:                       checkInterval,
		}
	case "grpc":
		check = GRPCCheck{
			DeregisterCriticalServiceAfter: deregisterInterval,
			GRPC:                           fmt.Sprintf("%s:%d", consulIp, port),
			Interval:                       checkInterval,
		}
	}

	register := ServiceRegisterInfo{
		Name:  name,
		Port:  port,
		ID:    fmt.Sprintf("%v-%v", name, port),
		Tags:  []string{tagGo},
		Check: check,
	}
	data, _ := json.Marshal(register)

	_, err := HttpPutWithHeader(url, headers, data)
	if err != nil {
		return err
	}

	return c.GetServicesToCache(name)
}

func (c *ConsulService) GetServices(serviceName string) ([]*ServiceInfo, error) {
	url := fmt.Sprintf("http://%s:%d/v1/health/service/%s", consulIp, consulPort, serviceName)
	headers := map[string]string{}

	result, err := HttpGetWithHeader(url, headers)
	if err != nil {
		return nil, err
	}

	var infos []ServiceHealthInfo
	if err := json.Unmarshal(result, &infos); err != nil {
		return nil, err
	}

	var validInfos []*ServiceInfo
	for _, info := range infos {
		for _, check := range info.Checks {
			if serviceName == check.ServiceName && check.Status == "passing" {
				validInfos = append(validInfos, &ServiceInfo{
					Address: info.Node.Address,
					Port:    info.Service.Port,
					Tags:    info.Service.Tags,
				})
			}
		}
	}

	return validInfos, nil
}

func (c *ConsulService) GetServicesToCache(serviceName string) error {
	infos, err := c.GetServices(serviceName)
	if err != nil {
		return err
	}

	return gConsulCache.Set(serviceName, infos)
}