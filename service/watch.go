package service

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func (s *Manager) WatchKey(key string, callback func(key string, value []byte)) error {
	var param = map[string]interface{}{
		"type":       "key",
		"key":        key,
		"datacenter": s.config.Datacenter,
		"token":      s.config.Token,
	}
	plan, _ := watch.Parse(param)
	plan.HybridHandler = func(blockParamVal watch.BlockingParamVal, val interface{}) {
		var v []byte
		if val != nil {
			v = (val).(*api.KVPair).Value
		}
		callback(key, v)
	}
	return plan.RunWithConfig(s.config.Address, s.config)
}

func (s *Manager) WatchService(service string, callback func(service string, serviceList *List)) error {
	var param = map[string]interface{}{
		"type":       "service",
		"service":    service,
		"datacenter": s.config.Datacenter,
		"token":      s.config.Token,
	}
	plan, _ := watch.Parse(param)
	plan.HybridHandler = func(blockParamVal watch.BlockingParamVal, val interface{}) {
		var v []*api.ServiceEntry
		if val != nil {
			v = (val).([]*api.ServiceEntry)
		}
		ss, _ := parseServices(service, v)
		callback(service, ss)
	}
	return plan.RunWithConfig(s.config.Address, s.config)
}
