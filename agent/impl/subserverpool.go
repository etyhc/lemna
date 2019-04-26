package impl

import (
	"lemna/agent"
	"lemna/agent/rpc"
	"lemna/config"
	configrpc "lemna/config/rpc"
	"lemna/logger"
)

// SubServerPool 无状态均衡服务器池，轮询方式调度服务器提供服务
type SubServerPool struct {
	servers    map[int32]map[string]*agent.Server
	clientPool agent.ClientPool
	scheduler  Scheduler
}

//Scheduler 服务器调度器
type Scheduler interface {
	schedule(map[string]*agent.Server, *agent.Client) *agent.Server
}

//GetTarget 目标池实现
func (ssp *SubServerPool) GetServer(target int32, client *agent.Client) *agent.Server {
	return ssp.scheduler.schedule(ssp.servers[target], client)
}

// SetTargetPool 目标池实现
func (ssp *SubServerPool) SetClientPool(cp agent.ClientPool) {
	ssp.clientPool = cp
}

func (ssp *SubServerPool) refreshServer(config *config.ServerConfig) bool {
	if servers, ok := ssp.servers[config.Type]; ok {
		if s, ok := servers[config.Addr]; ok {
			s.Info = config
			return ok
		}
	}
	return false
}

func (ssp *SubServerPool) registerServer(config *config.ServerConfig) {
	go func() {
		cs := &rpc.ClientService{Addr: config.Addr, Typeid: config.Type}
		err := cs.Init()
		if err == nil {
			if _, ok := ssp.servers[cs.Typeid]; !ok {
				ssp.servers[cs.Typeid] = make(map[string]*agent.Server)
			}
			s := agent.NewServer(cs.Forwarder(), cs.TypeID(), config)
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
func (ssp *SubServerPool) SubscribeServer(addr string, sch Scheduler) error {
	ssp.servers = make(map[int32]map[string]*agent.Server)
	ssp.scheduler = sch
	finder := configrpc.ChannelUser{Addr: addr}
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
