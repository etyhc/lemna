package impl

import (
	"lemna/agent"
	"lemna/logger"
	"reflect"
)

// BalanceScheduler 调度器，轮流方式调度服务器
type BalanceScheduler struct {
	last int
}

func (bs *BalanceScheduler) schedule(servers map[string]*agent.Server, client *agent.Client) *agent.Server {
	if len(servers) > 0 {
		keys := reflect.ValueOf(servers).MapKeys()
		return servers[keys[bs.algorithm(len(servers))].String()]
	}
	return nil
}

//algorithm 轮询模式选择
//BUG(yhc): 不是严格轮询
func (bs *BalanceScheduler) algorithm(all int) int {
	if all == 0 {
		return 0
	}
	bs.last = (bs.last + 1) % all
	logger.Debugf("all%d select %d", all, bs.last)
	return bs.last
}
