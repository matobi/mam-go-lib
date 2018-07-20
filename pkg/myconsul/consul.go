package myconsul

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type ServiceAddress struct {
	Name     string `json:"Name"`
	IP       string `json:"IP"`
	Port     int    `json:"Port"`
	Location string `json:"Location"`
}

type consulService struct {
	Service string `json:"Service"`
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

func ListServices(client *http.Client, consulAddress string) (map[string]ServiceAddress, error) {
	url := fmt.Sprintf("%s/v1/agent/services", consulAddress)
	resp, err := client.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed call url; url=%s", url)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.Errorf("failed call consul; url=%s; code=%d", url, resp.StatusCode)
	}
	consulServices := make(map[string]consulService)
	if err := json.NewDecoder(resp.Body).Decode(&consulServices); err != nil {
		return nil, errors.Wrapf(err, "Failed read response; url=%s", url)
	}
	services := make(map[string]ServiceAddress)
	for _, v := range consulServices {
		services[v.Service] = ServiceAddress{
			Name:     v.Service,
			IP:       v.Address,
			Port:     v.Port,
			Location: fmt.Sprintf("%s:%d", v.Address, v.Port),
		}
	}
	return services, nil
}

func FindService(h *http.Client, consulAddress string, serviceName string) (*ServiceAddress, error) {
	services, err := ListServices(h, consulAddress)
	if err != nil {
		return nil, err
	}
	svc, isFound := services[serviceName]
	if !isFound {
		return nil, errors.Errorf("service not in consul; service=%s; consul=%s", serviceName, consulAddress)
	}
	return &svc, nil
}
