// Package server 实现了服务器的rpc调用.
package server

import (
	"fmt"
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
	servers map[int32]map[string]*Server //所有订阅并连接的服务器
	cp      agent.TargetPool             //客户端池
	addr    string                       //订阅服务器地址
}

// schedule 根据服务器信息调度服务器
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

// NewSubServerPool 新订阅服务器池
func NewSubServerPool(addr string) *SubServerPool {
	return &SubServerPool{servers: make(map[int32]map[string]*Server), addr: addr}
}

// GetTarget 服务器池接口实现
func (ssp *SubServerPool) GetTarget(target int32) agent.Target {
	ret := schedule(ssp.servers[target])
	if ret != nil {
		return ret
	}
	return nil
}

// Bind 服务器池接口实现
func (ssp *SubServerPool) Bind(cp agent.TargetPool) {
	ssp.cp = cp
}

// refreshServer 刷新服务器信息
//               如果已经有服务器，更新服务器信息，返回true，否则返回false
func (ssp *SubServerPool) refreshServer(info *Info) bool {
	if servers, ok := ssp.servers[info.Type]; ok {
		if s, ok := servers[info.Addr]; ok {
			s.Info = info
			return ok
		}
	}
	return false
}

// registerServer 注册服务器
//                根据服务器信息，连接服务器,并待用服务器转发函数
func (ssp *SubServerPool) registerServer(info *Info) {
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
			err = agent.S2C(s, ssp.cp)
			delete(ssp.servers[info.Type], info.Addr)
		}
		logger.Error(err)
	}()
}

// Run 从配置频道服务器订阅服务器信息
//     订阅到服务器后链接注册服务器，如果已注册更新服务器信息
func (ssp *SubServerPool) Run() error {
	finder := crpc.Channel{Addr: ssp.addr}
	ch, err := finder.Subscribe(&Info{})
	if err != nil {
		return err
	}
	for {
		info, ok := (<-ch).(*Info)
		if !ok {
			return fmt.Errorf("subscribe closed")
		}
		if ssp.refreshServer(info) {
			logger.Info("refresh ", info)
		} else {
			logger.Info("register ", info)
			ssp.registerServer(info)
		}
	}
}
