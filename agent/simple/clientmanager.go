package simple

import (
	fmt "fmt"
	"lemna/agent"
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

func (cm *SimpleClientManager) NewClient(sessionid int32) (c agent.Client, err error) {
	var ok bool
	cm.rw.Lock()
	c, ok = cm.clients[sessionid]
	if !ok {
		c = &SimpleClient{sessionid: sessionid}
		cm.clients[sessionid] = c
	}
	cm.rw.Unlock()
	return
}

func (cm *SimpleClientManager) GetClient(sessionid int32) (c agent.Client, err error) {
	var ok bool
	cm.rw.RLock()
	if c, ok = cm.clients[sessionid]; !ok {
		err = fmt.Errorf("not find client sessionid=%d", sessionid)
	}
	cm.rw.RUnlock()
	return
}

func (cm *SimpleClientManager) DelClient(sessionid int32) {
	cm.rw.Lock()
	delete(cm.clients, sessionid)
	cm.rw.Unlock()

}
