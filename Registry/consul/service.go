package consul

import (
	"encoding/json"
	"fmt"
)

const (
	baseUrl  = "http://127.0.0.1:8500"
	localUrl = "127.0.0.1"
)

const (
	HEALTH_NIL  = 0
	HEALTH_HTTP = 1
	HEALTH_GRPC = 2
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

func (this *ConsulService) Register(name string, port int64, healthType string) error {
	url := fmt.Sprintf("%s/v1/agent/service/register", baseUrl)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	var check interface{}
	switch healthType {
	case "http":
		check = HTTPCheck{
			DeregisterCriticalServiceAfter: "10m",
			HTTP:                           fmt.Sprintf("http://%s:%d/health", localUrl, port),
			Interval:                       "30s",
		}
	case "grpc":
		check = GRPCCheck{
			DeregisterCriticalServiceAfter: "10m",
			GRPC:                           fmt.Sprintf("%s:%d", localUrl, port),
			Interval:                       "30s",
		}
	}

	register := ServiceRegisterInfo{
		Name:  name,
		Port:  port,
		ID:    fmt.Sprintf("%v-%v", name, port),
		Tags:  []string{"go"},
		Check: check,
	}
	data, _ := json.Marshal(register)

	_, err := HttpPutWithHeader(url, headers, data)
	if err != nil {
		return err
	}

	return this.GetServicesToCache(name)
}

func (this *ConsulService) GetServices(serviceName string) ([]*ServiceInfo, error) {
	url := fmt.Sprintf("%s/v1/health/service/%s", baseUrl, serviceName)
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

func (this *ConsulService) GetServicesToCache(serviceName string) error {
	infos, err := this.GetServices(serviceName)
	if err != nil {
		return err
	}

	return g_consulCache.Set(serviceName, infos)
}
