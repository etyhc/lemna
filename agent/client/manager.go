package client

import (
	fmt "fmt"
	"lemna/arpc"
	"sync"
)

//manager 仅供client.Service使用
type manager struct {
	clients map[uint32]*Target
	mu      sync.Mutex
}

func newMananger() *manager {
	return &manager{clients: make(map[uint32]*Target)}
}

func (cm *manager) newTarget(s arpc.CAgent_ForwardServer, id uint32) (*Target, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.clients[id]
	if !ok {
		cm.clients[id] = newTarget(s, id)
		return cm.clients[id], nil
	}
	return nil, fmt.Errorf("repeated client<id=%d>", id)
}

func (cm *manager) getTarget(id uint32) *Target {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	ret, ok := cm.clients[id]
	if ok {
		return ret
	}
	return nil
}

func (cm *manager) delTarget(id uint32) {
	cm.mu.Lock()
	delete(cm.clients, id)
	cm.mu.Unlock()
}
