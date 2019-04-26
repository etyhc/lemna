package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"sync"
)

type clientManager struct {
	clients    map[int32]*Client
	serverPool ServerPool
	mu         sync.Mutex
}

func newClientMananger() *clientManager {
	return &clientManager{clients: make(map[int32]*Client)}
}

// GetClient 目标池接口实现
func (cm *clientManager) GetClient(cid int32, s *Server) *Client {
	return cm.getClient(cid)
}

// SetServerPool 目标池接口实现
func (cm *clientManager) SetServerPool(sp ServerPool) {
	cm.serverPool = sp
}

func (cm *clientManager) newClient(s rpc.Client_ForwardServer, id int32) (*Client, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.clients[id]
	if !ok {
		cm.clients[id] = NewClient(s, id)
		return cm.clients[id], nil
	}
	return nil, fmt.Errorf("repeated sessionid")
}

func (cm *clientManager) getClient(id int32) *Client {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if ret, ok := cm.clients[id]; ok {
		return ret
	}
	return nil
}

func (cm *clientManager) delClient(id int32) {
	cm.mu.Lock()
	delete(cm.clients, id)
	cm.mu.Unlock()
}
