package impl

import (
	"lemna/agent"
	"lemna/agent/rpc"
	"lemna/config"
	configrpc "lemna/config/rpc"
	"lemna/logger"
)

type server struct {
	target *agent.Target
	info   *config.ServerConfig
}

// ServerPool 无状态均衡服务器池，轮询方式调度服务器提供服务
type ServerPool struct {
	servers    map[int32]map[string]*server
	clientPool agent.TargetPool
	scheduler  Scheduler
}

//Scheduler 服务器调度器
type Scheduler interface {
	schedule(map[string]*server, *agent.Target) *agent.Target
}

//GetTarget 目标池实现
func (sp *ServerPool) GetTarget(target int32, client *agent.Target) *agent.Target {
	return sp.scheduler.schedule(sp.servers[target], client)
}

// SetTargetPool 目标池实现
func (sp *ServerPool) SetTargetPool(tp agent.TargetPool) {
	sp.clientPool = tp
}

func (sp *ServerPool) refreshServer(config *config.ServerConfig) bool {
	if servers, ok := sp.servers[config.Type]; ok {
		s, ok := servers[config.Addr]
		s.info = config
		return ok
	}
	return false
}

func (sp *ServerPool) registerServer(config *config.ServerConfig) {
	go func() {
		cs := &rpc.ClientService{Addr: config.Addr, Typeid: config.Type}
		err := cs.Init()
		if err == nil {
			if _, ok := sp.servers[cs.Typeid]; !ok {
				sp.servers[cs.Typeid] = make(map[string]*server)
			}
			s := &server{agent.NewTarget(cs.Stream(), cs.TypeID()), config}
			sp.servers[cs.Typeid][cs.Addr] = s
			logger.Infof("%s(%d) is running", cs.Addr, cs.Typeid)
			err = s.target.Run(sp.clientPool)
			delete(sp.servers[cs.Typeid], cs.Addr)
		}
		logger.Error(err)
	}()
}

//SubscribeServer 从配置频道服务器订阅服务器信息
//                订阅到服务器后链接注册服务器，如果已注册更新服务器信息
func (sp *ServerPool) SubscribeServer(addr string, sch Scheduler) error {
	sp.servers = make(map[int32]map[string]*server)
	sp.scheduler = sch
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
			if sp.refreshServer(info) {
				logger.Info("refresh ", info)
			} else {
				logger.Info("register ", info)
				sp.registerServer(info)
			}
		}
	}()
	return nil
}
