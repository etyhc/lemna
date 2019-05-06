package server

import (
	"lemna/agent"
	"lemna/agent/arpc"
	"lemna/content/crpc"
	"lemna/logger"
)

//SubServerPool 订阅服务器信息并连接的服务器池
//              根据服务器调度信息，提供服务器调度，被调度的服务器都应该是无状态的
//              如果是有状态服务器，那么服务器应设置成SERVERSCHENIL,不让代理服务器进行调度
//              并自己实现有状态调度器，与客户端协调状态
type SubServerPool struct {
	servers map[int32]map[string]*Server
	cp      agent.TargetPool
	addr    string
}

func schedule(servers map[string]*Server) *Server {
	var ret *Server
	for _, server := range servers {
		switch server.Info.Sche {
		case SERVERSCHEROUND:
			if ret == nil || ret.Round > server.Round {
				ret = server
			}
		case SERVERSCHELOAD:
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

func NewSubServerPool(addr string) *SubServerPool {
	return &SubServerPool{servers: make(map[int32]map[string]*Server), addr: addr}
}

//GetServer 服务器池接口实现
func (ssp *SubServerPool) GetTarget(target int32) agent.Target {
	ret := schedule(ssp.servers[target])
	if ret != nil {
		return ret
	}
	return nil
}

//SetClientPool 服务器池接口实现
func (ssp *SubServerPool) Bind(cp agent.TargetPool) {
	ssp.cp = cp
}

func (ssp *SubServerPool) refreshServer(info *ServerInfo) bool {
	if servers, ok := ssp.servers[info.Type]; ok {
		if s, ok := servers[info.Addr]; ok {
			s.Info = info
			return ok
		}
	}
	return false
}

func (ssp *SubServerPool) registerServer(info *ServerInfo) {
	go func() {
		c := arpc.NewClient(info.Addr, info.Type, agent.ServiceID)
		err := c.Init()
		if err == nil {
			if _, ok := ssp.servers[c.TypeID()]; !ok {
				ssp.servers[c.TypeID()] = make(map[string]*Server)
			}
			s := NewServer(c, info)
			ssp.servers[info.Type][info.Addr] = s
			logger.Infof("%s(%d) is running", info.Addr, info.Type)
			err = agent.StoC(s, ssp.cp)
			delete(ssp.servers[info.Type], info.Addr)
		}
		logger.Error(err)
	}()
}

//Start 从配置频道服务器订阅服务器信息
//      订阅到服务器后链接注册服务器，如果已注册更新服务器信息
func (ssp *SubServerPool) Run() error {
	finder := crpc.Channel{Addr: ssp.addr}
	ch, err := finder.Subscribe(&ServerInfo{})
	if err != nil {
		return err
	}
	go func() {
		for {
			info, ok := (<-ch).(*ServerInfo)
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
