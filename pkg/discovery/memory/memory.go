package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"pinterest2.0/pkg/registry"
)

type serviceName string
type instanceID string

type MemoryRegistry struct {
	sync.RWMutex
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *MemoryRegistry {
	return &MemoryRegistry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *MemoryRegistry) Register(ctx context.Context, instanceId string, serviceN string, hostPort string) error {
	r.Lock()
	defer r.Unlock()
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)
	if _, ok := r.serviceAddrs[sName]; !ok {
		r.serviceAddrs[sName] = make(map[instanceID]*serviceInstance)
	}
	r.serviceAddrs[sName][iID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()}
	return nil
}

func (r *MemoryRegistry) Deregister(ctx context.Context, instanceId string, serviceN string) error {
	r.Lock()
	defer r.Unlock()
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)
	if _, ok := r.serviceAddrs[sName]; !ok {
		return nil
	}
	delete(r.serviceAddrs[sName], iID)
	return nil
}

func (r *MemoryRegistry) ReportHealthyState(instanceId string, serviceN string) error {
	r.Lock()
	defer r.Unlock()
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)
	if _, ok := r.serviceAddrs[sName]; !ok {
		return errors.New("Service is not registered yet")
	}
	if _, ok := r.serviceAddrs[sName][iID]; !ok {
		return errors.New("Service instance is not registered yet")
	}
	r.serviceAddrs[sName][iID].lastActive = time.Now()
	return nil
}

func (r *MemoryRegistry) ServiceAddress(ctx context.Context, serviceN string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()
	sName := serviceName(serviceN)
	if len(r.serviceAddrs[sName]) == 0 {
		return nil, registry.ErrNotFound
	}
	var res []string
	for _, i := range r.serviceAddrs[sName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		res = append(res, i.hostPort)
	}
	return res, nil
}
