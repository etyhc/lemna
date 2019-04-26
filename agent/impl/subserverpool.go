package impl

import (
	"lemna/agent"
	"lemna/agent/rpc"
	"lemna/config"
	configrpc "lemna/config/rpc"
	"lemna/logger"
)

//SubServerPool 订阅服务器信息并连接的服务器池
//              根据服务器调度信息，提供服务器调度，被调度的服务器都应该是无状态的
//              如果是有状态服务器，那么服务器应设置成SERVERSCHENIL,不让代理服务器进行调度
//              并自己实现有状态调度器，与客户端协调状态
type SubServerPool struct {
	servers    map[int32]map[string]*agent.Server
	clientPool agent.ClientPool
}

func schedule(servers map[string]*agent.Server, client *agent.Client) *agent.Server {
	var ret *agent.Server
	for _, server := range servers {
		switch server.Info.Sche {
		case config.SERVERSCHEROUND:
			if ret == nil || ret.Round > server.Round {
				ret = server
			}
		case config.SERVERSCHELOAD:
			if ret == nil || ret.Info.Load > server.Info.Load {
				ret = server
			}
		}
	}
	if ret != nil {
		ret.Round++
	}
	return ret
}

//GetServer 服务器池接口实现
func (ssp *SubServerPool) GetServer(target int32, client *agent.Client) *agent.Server {
	return schedule(ssp.servers[target], client)
}

//SetClientPool 服务器池接口实现
func (ssp *SubServerPool) SetClientPool(cp agent.ClientPool) {
	ssp.clientPool = cp
}

func (ssp *SubServerPool) refreshServer(info *config.ServerInfo) bool {
	if servers, ok := ssp.servers[info.Type]; ok {
		if s, ok := servers[info.Addr]; ok {
			s.Info = info
			return ok
		}
	}
	return false
}

func (ssp *SubServerPool) registerServer(info *config.ServerInfo) {
	go func() {
		cs := &rpc.ClientService{Addr: info.Addr, Typeid: info.Type}
		err := cs.Init()
		if err == nil {
			if _, ok := ssp.servers[cs.Typeid]; !ok {
				ssp.servers[cs.Typeid] = make(map[string]*agent.Server)
			}
			s := agent.NewServer(cs.Forwarder(), cs.TypeID(), info)
			ssp.servers[cs.Typeid][cs.Addr] = s
			logger.Infof("%s(%d) is running", cs.Addr, cs.Typeid)
			err = s.Run(ssp.clientPool)
			delete(ssp.servers[cs.Typeid], cs.Addr)
		}
		logger.Error(err)
	}()
}

//SubscribeServer 从配置频道服务器订阅服务器信息
//                订阅到服务器后链接注册服务器，如果已注册更新服务器信息
func (ssp *SubServerPool) SubscribeServer(addr string) error {
	ssp.servers = make(map[int32]map[string]*agent.Server)
	finder := configrpc.ChannelUser{Addr: addr}
	ch, err := finder.Subscribe("server", &config.ServerInfo{})
	if err != nil {
		return err
	}
	go func() {
		for {
			info, ok := (<-ch).(*config.ServerInfo)
			if !ok {
				logger.Debug("subscribe closed")
				return
			}
			if ssp.refreshServer(info) {
				logger.Info("refresh ", info)
			} else {
				logger.Info("register ", info)
				ssp.registerServer(info)
			}
		}
	}()
	return nil
}
