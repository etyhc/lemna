package simple

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/agent/rpc"
	"sync"
)

type SimpleClientManager struct {
	clients map[int32]agent.Client
	rw      sync.RWMutex
}

func NewSimpleClientManager() (cm *SimpleClientManager) {
	cm = &SimpleClientManager{}
	cm.clients = make(map[int32]agent.Client)
	return
}

func (cm *SimpleClientManager) NewClient(sessionid int32) (agent.Client, error) {
	cm.rw.Lock()
	c, ok := cm.clients[sessionid]
	if !ok {
		c = &SimpleClient{sessionid: sessionid, servers: make(map[int32]*rpc.ClientService)}
		cm.clients[sessionid] = c
	}
	cm.rw.Unlock()
	return c, nil
}

func (cm *SimpleClientManager) GetClient(sessionid int32) (agent.Client, error) {
	cm.rw.RLock()
	c, ok := cm.clients[sessionid]
	cm.rw.RUnlock()
	if ok {
		return c, nil
	}
	return c, fmt.Errorf("not find client sessionid=%d", sessionid)
}

func (cm *SimpleClientManager) DelClient(sessionid int32) {
	cm.rw.Lock()
	delete(cm.clients, sessionid)
	cm.rw.Unlock()

}
