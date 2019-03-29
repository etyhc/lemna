package main

import (
	"lemna/agent"
	"sync"
)

type clientManager struct {
	clients map[int32]agent.Client
	rw      sync.RWMutex
}

func (cm *clientManager) GetClient(sessionid int32) (c agent.Client, err error) {
	var ok bool
	cm.rw.RLock()
	c, ok = cm.clients[sessionid]
	cm.rw.RUnlock()
	if !ok {
		c = &client{sessionid: sessionid}
		cm.rw.Lock()
		cm.clients[sessionid] = c
		cm.rw.Unlock()
	}
	return
}

func (cm *clientManager) DelClient(sessionid int32) {
	cm.rw.Lock()
	delete(cm.clients, sessionid)
	cm.rw.Unlock()
}
