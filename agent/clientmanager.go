package agent

import "sync"

type clientManager struct {
	clients    map[int32]*Target
	serverPool TargetPool
	cmu        sync.Mutex
}

func newClientMananger() *clientManager {
	return &clientManager{clients: make(map[int32]*Target)}
}

// GetTarget 目标池接口实现
func (tm *clientManager) GetTarget(tt int32, t *Target) *Target {
	return tm.getClient(tt)
}

// SetTargetPool 目标池接口实现
func (tm *clientManager) SetTargetPool(tp TargetPool) {
	tm.serverPool = tp
}

func (tm *clientManager) newClient(s Stream, id int32) *Target {
	tm.cmu.Lock()
	defer tm.cmu.Unlock()
	tm.clients[id] = NewTarget(s, id)
	return tm.clients[id]
}

func (tm *clientManager) getClient(id int32) *Target {
	tm.cmu.Lock()
	defer tm.cmu.Unlock()
	if ret, ok := tm.clients[id]; ok {
		return ret
	}
	return nil
}

func (tm *clientManager) delClient(id int32) {
	tm.cmu.Lock()
	delete(tm.clients, id)
	tm.cmu.Unlock()
}
