package impl

import (
	"lemna/agent"
	"reflect"
)

type BalanceScheduler struct {
	last int
}

func (bs *BalanceScheduler) schedule(servers map[string]*server, client *agent.Target) *agent.Target {
	if len(servers) > 0 {
		keys := reflect.ValueOf(servers).MapKeys()
		return servers[keys[bs.algorithm(len(servers))].String()].target
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
	return bs.last
}
