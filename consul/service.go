package consul

import (
	"encoding/json"
	"fmt"
)

const (
	baseUrl  = "http://127.0.0.1:8500"
	localUrl = "http://127.0.0.1"
)

type ConsulService struct {
}

type ServiceRegistration struct {
	ID    string        `json:"ID"`
	Name  string        `json:"Name"`
	Port  int64         `json:"Port"`
	Tags  []string      `json:"Tags"`
	Check RegisterCheck `json:"Check"`
}

type RegisterCheck struct {
	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	HTTP                           string `json:"HTTP"`
	Interval                       string `json:"Interval"`
}

type ServiceInfo struct {
	ID      string   `json:"ID"`
	Name    string   `json:"Name"`
	Address string   `json:"Address"`
	Port    int64    `json:"Port"`
	Tags    []string `json:"Tags"`
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

func NewConsulService() *ConsulService {
	return &ConsulService{}
}

func (this *ConsulService) Register(serviceName string, port int64) error {
	url := fmt.Sprintf("%s/v1/agent/service/register", baseUrl)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	register := ServiceRegistration{
		Name: serviceName,
		Port: port,
		ID:   fmt.Sprintf("%v-%v", serviceName, port),
		Tags: []string{"go"},
		Check: RegisterCheck{
			DeregisterCriticalServiceAfter: "10m",
			HTTP:                           fmt.Sprintf("%s:%d/health", localUrl, port),
			Interval:                       "30s",
		},
	}
	data, _ := json.Marshal(register)

	_, err := HttpPutWithHeader(url, headers, data)
	if err != nil {
		return err
	}

	infos, err := this.GetServices(serviceName)
	if err != nil {
		return err
	}

	return g_consulCache.Set(serviceName, infos...)
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
					ID:      info.Service.ID,
					Name:    info.Service.Name,
					Address: info.Node.Address,
					Port:    info.Service.Port,
					Tags:    info.Service.Tags,
				})
			}
		}
	}
	return validInfos, err
}
