package client

import (
	fmt "fmt"
	"lemna/agent/arpc"
	"sync"
)

type clientManager struct {
	clients map[int32]*Client
	mu      sync.Mutex
}

func newClientMananger() *clientManager {
	return &clientManager{clients: make(map[int32]*Client)}
}

func (cm *clientManager) newClient(s arpc.Client_ForwardServer, id int32) (*Client, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.clients[id]
	if !ok {
		cm.clients[id] = NewClient(s, id)
		return cm.clients[id], nil
	}
	return nil, fmt.Errorf("repeated client<id=%d>", id)
}

func (cm *clientManager) getClient(id int32) *Client {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	ret, ok := cm.clients[id]
	if ok {
		return ret
	}
	return nil
}

func (cm *clientManager) delClient(id int32) {
	cm.mu.Lock()
	delete(cm.clients, id)
	cm.mu.Unlock()
}
