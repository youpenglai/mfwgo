package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	consulIp   = "127.0.0.1"
	consulPort = 8500
)

// 支持从环境变量中获取
func initConsul() {
	ip := os.Getenv("MFW_CONSUL_IP")
	if ip != "" {
		consulIp = ip
	}
	port := os.Getenv("MFW_CONSUL_PORT")
	if port != "" {
		v, e := strconv.ParseInt(port, 10, 32)
		if e != nil {
			return
		}
		consulPort = int(v)
	}
}

const (
	deregisterInterval = "10m"
	tagGo              = "go"
	checkInterval      = "30s"
)

type ServiceRegisterInfo struct {
	ID                string      `json:"ID"`
	Name              string      `json:"Name"`
	Port              int64       `json:"Port"`
	Tags              []string    `json:"Tags"`
	Check             interface{} `json:"Check"`
	EnableTagOverride bool        `json:"EnableTagOverride"`
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
		Name:              name,
		Port:              port,
		ID:                fmt.Sprintf("%v-%v", name, port),
		Tags:              []string{tagGo},
		Check:             check,
		EnableTagOverride: true,
	}
	data, _ := json.Marshal(register)

	errMsg, err := HttpPutWithHeader(url, headers, data)
	if err != nil {
		return err
	}

	if string(errMsg) != "" {
		return errors.New(string(errMsg))
	}

	return c.GetServicesToCache(name)
}

func (c *ConsulService) GetServices(serviceName string) ([]*ServiceInfo, error) {
	url := fmt.Sprintf("http://%s:%d/v1/health/service/%s?passing", consulIp, consulPort, serviceName)
	headers := map[string]string{}

	result, err := HttpGetWithHeader(url, headers)
	if err != nil {
		return nil, err
	}

	var infos []ServiceHealthInfo
	if err := json.Unmarshal(result, &infos); err != nil {
		return nil, err
	}

	//if len(infos) == 0 {
	//	return nil, errors.New("not found service: " + serviceName)
	//}

	var validInfos []*ServiceInfo
	for _, info := range infos {
		validInfos = append(validInfos, &ServiceInfo{
			Address: info.Node.Address,
			Port:    info.Service.Port,
			Tags:    info.Service.Tags,
		})
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

func init() {
	initConsul()
}
