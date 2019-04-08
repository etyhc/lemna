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

type SimpleBalancer struct {
	servers map[int32]map[string]*rpc.ClientService
	as      *agent.Service
}

func (sb *SimpleBalancer) GetServer(target int32, c agent.Client) (agent.Server, bool) {
	r := rand.New(rand.NewSource(99))
	ss, ok := sb.servers[target]
	if ok && len(ss) > 0 {
		keys := reflect.ValueOf(ss).MapKeys()
		return ss[keys[r.Intn(len(ss))].String()], true
	}
	return nil, false
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
	sb.servers = make(map[int32]map[string]*rpc.ClientService)
	return
}

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
			logger.Info("register ", info)
			sb.registerServer(&rpc.ClientService{Typeid: info.Type, Addr: info.Addr})
		}
	}()
	return nil
}

func (sb *SimpleBalancer) Start(as *agent.Service) {
	sb.as = as
	sb.subscribe()
}
