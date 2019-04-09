package simple

import (
	"lemna/agent"
	"lemna/agent/rpc"
	"lemna/config"
	configrpc "lemna/config/rpc"
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

// GetServer 如果客户端没有相应类型服务器,根据算法得到类型为target的服务器
func (sb *SimpleBalancer) GetServer(target int32, c agent.Client) (agent.Server, bool) {
	if server, ok := c.(*SimpleClient).servers[target]; ok {
		select {
		case <-server.Stream().Context().Done(): //保存的服务器失效
			delete(c.(*SimpleClient).servers, target)
		default:
			return server, ok
		}
	}
	if ss, ok := sb.servers[target]; ok && len(ss) > 0 {
		keys := reflect.ValueOf(ss).MapKeys()
		c.(*SimpleClient).servers[target] = ss[keys[sb.algorithm(len(ss))].String()]
		return c.(*SimpleClient).servers[target], true
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
	finder := configrpc.ChannelUser{Addr: configrpc.ConfigServerAddr}
	ch, err := finder.Subscribe("server", &config.ServerConfig{})
	if err != nil {
		return err
	}
	go func() {
		for {
			info, ok := (<-ch).(*config.ServerConfig)
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
	err := sb.subscribe()
	if err != nil {
		logger.Error(err)
	}
}
