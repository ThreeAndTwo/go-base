package service

import (
	"github.com/hashicorp/consul/api"
	"os"
	"strconv"
	"time"
)

type Info struct {
	ServiceName string
	ServiceID   string
	IP          string
	Port        int
	Timestamp   int
}

type List []Info

type Config struct {
	Addr       string `mapstructure:"addr"`
	Token      string `mapstructure:"token"`
	Datacenter string `mapstructure:"data_center"`
	WatchKey   string `mapstructure:"watch_key"`
}

type Manager struct {
	config      *api.Config
	serviceInfo *Info
	client      *api.Client
	quitChan    chan os.Signal
	checker     *api.AgentServiceCheck
}

func NewManager(config *Config) *Manager {
	_config := api.DefaultConfig()
	_config.Address = config.Addr
	_config.Token = config.Token
	if config.Datacenter != "" {
		_config.Datacenter = config.Datacenter
	}
	_client, err := api.NewClient(_config)
	if err != nil {
		return nil
	}
	return &Manager{config: _config, client: _client, quitChan: make(chan os.Signal, 1)}
}

func (s *Manager) SetServiceInfo(serviceInfo *Info) {
	s.serviceInfo = serviceInfo
}

func (s *Manager) getServiceId() string {
	return s.serviceInfo.ServiceName + "-" + s.serviceInfo.IP + "-" + strconv.Itoa(s.serviceInfo.Port)
}

func (s *Manager) AddChecker(checker *api.AgentServiceCheck) {
	s.checker = checker
}

func (s *Manager) Start() error {
	var tags []string
	service := &api.AgentServiceRegistration{
		ID:      s.getServiceId(),
		Name:    s.serviceInfo.ServiceName,
		Port:    s.serviceInfo.Port,
		Address: s.serviceInfo.IP,
		Tags:    tags,
		Check:   s.checker,
	}
	err := s.client.Agent().ServiceRegister(service)
	if err != nil {
		return err
	}
	go s.WaitToUnDeregisterService()
	return nil
}

func (s *Manager) Close() {
	s.quitChan <- os.Kill
}

func (s *Manager) WaitToUnDeregisterService() {
	<-s.quitChan
	if s.client == nil {
		return
	}
	if err := s.client.Agent().ServiceDeregister(s.getServiceId()); err != nil {
		println("deregister service failed, " + err.Error())
	}
}

func (s *Manager) Discover(serviceName string) (*List, error) {
	servicesData, _, err := s.client.Health().Service(serviceName, "", true,
		&api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	return parseServices(serviceName, servicesData)
}

func parseServices(serviceName string, servicesData []*api.ServiceEntry) (*List, error) {
	var serviceList List
	for _, entry := range servicesData {
		if serviceName != entry.Service.Service {
			continue
		}
		for _, health := range entry.Checks {
			if health.ServiceName != serviceName || health.Status != "passing" {
				continue
			}
			var node = Info{
				ServiceName: entry.Service.Service,
				ServiceID:   entry.Service.ID,
				IP:          entry.Service.Address,
				Port:        entry.Service.Port,
				Timestamp:   time.Now().Nanosecond(),
			}
			serviceList = append(serviceList, node)
		}
	}
	return &serviceList, nil
}
