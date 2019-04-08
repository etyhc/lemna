package simple

import (
	"lemna/agent"
	"lemna/agent/rpc"
	configrpc "lemna/config/rpc"
	"lemna/config/rpc/server"
	"lemna/logger"
	"math/rand"
	"reflect"
)

// SimpleBalancer 简单的无状态均衡器，轮询服务器提供服务
type SimpleBalancer struct {
	servers map[int32]map[string]*rpc.ClientService
	as      *agent.Service
	rand    *rand.Rand
}

var last = 0

//轮询模式选择,TODO有bug
func (sb *SimpleBalancer) algorithm(all int) int {
	if all == 0 {
		return 0
	}
	last = (last + 1) % all
	return last
}

// GetServer 根据算法得到类型为target的服务器,c无视了
func (sb *SimpleBalancer) GetServer(target int32, c agent.Client) (agent.Server, bool) {
	ss, ok := sb.servers[target]
	if ok && len(ss) > 0 {
		keys := reflect.ValueOf(ss).MapKeys()
		return ss[keys[sb.algorithm(len(ss))].String()], true
	}
	return nil, false
}

func (sb *SimpleBalancer) hasServer(typeid int32, addr string) bool {
	if servers, ok := sb.servers[typeid]; ok {
		_, ok := servers[addr]
		return ok
	}
	return false
}

func (sb *SimpleBalancer) registerServer(cs *rpc.ClientService) {
	go func() {
		err := cs.Init()
		if err == nil {
			if _, ok := sb.servers[cs.Typeid]; !ok {
				sb.servers[cs.Typeid] = make(map[string]*rpc.ClientService)
			}
			sb.servers[cs.Typeid][cs.Addr] = cs
			logger.Infof("%s(%d) is running", cs.Addr, cs.Typeid)
			err = sb.as.RunServer(cs)
			delete(sb.servers[cs.Typeid], cs.Addr)
		}
		logger.Error(err)
	}()
}

func NewSimpleBalancer() (sb *SimpleBalancer) {
	sb = &SimpleBalancer{}
	sb.rand = rand.New(rand.NewSource(99))
	sb.servers = make(map[int32]map[string]*rpc.ClientService)
	return
}

//从配置频道服务器得到服务器信息
func (sb *SimpleBalancer) subscribe() error {
	finder := configrpc.ChannelClient{Addr: configrpc.ConfigServerAddr}
	ch, err := finder.Subscribe("server", &server.Config{})
	if err != nil {
		return err
	}
	go func() {
		for {
			info, ok := (<-ch).(*server.Config)
			if !ok {
				logger.Debug("subscribe closed")
				return
			}
			if !sb.hasServer(info.Type, info.Addr) {
				logger.Info("register ", info)
				sb.registerServer(&rpc.ClientService{Typeid: info.Type, Addr: info.Addr})
			}
		}
	}()
	return nil
}

func (sb *SimpleBalancer) Start(as *agent.Service) {
	sb.as = as
	sb.subscribe()
}
