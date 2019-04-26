package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"sync"
	"time"
)

type clientManager struct {
	clients    map[int32]*Client
	serverPool ServerPool
	mu         sync.Mutex
}

func newClientMananger() *clientManager {
	ret := &clientManager{clients: make(map[int32]*Client)}
	go func() {
		t := time.NewTicker(time.Duration(time.Second * 5))
		for {
			<-t.C
			ret.mu.Lock()
			for _, client := range ret.clients {
				if client.stream == nil {
					delete(ret.clients, client.id)
				}
			}
			ret.mu.Unlock()
		}
	}()
	return ret
}

// GetClient 目标池接口实现
func (cm *clientManager) GetClient(cid int32, s *Server) *Client {
	return cm.getClient(cid)
}

// SetServerPool 目标池接口实现
func (cm *clientManager) SetServerPool(sp ServerPool) {
	cm.serverPool = sp
}

func (cm *clientManager) newClient(id int32) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	_, ok := cm.clients[id]
	if !ok {
		cm.clients[id] = NewClient(nil, id)
	}
	return !ok
}

func (cm *clientManager) initClient(s rpc.Client_ForwardServer, id int32) (*Client, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	client, ok := cm.clients[id]
	if ok {
		if client.stream == nil {
			client.stream = s
			return cm.clients[id], nil
		}
		return nil, fmt.Errorf("init error:repeated init client<id=%d>", id)
	}
	return nil, fmt.Errorf("init error:client<id=%d> not login", id)
}

func (cm *clientManager) getClient(id int32) *Client {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	ret, ok := cm.clients[id]
	if ok && ret.stream != nil {
		return ret
	}
	return nil
}

func (cm *clientManager) delClient(id int32) {
	cm.mu.Lock()
	delete(cm.clients, id)
	cm.mu.Unlock()
}
