package impl

import (
	"lemna/agent"
	"lemna/agent/rpc"
	"lemna/config"
	configrpc "lemna/config/rpc"
	"lemna/logger"
	"math/rand"
	"reflect"
	"sync"
)

// BalancePool 无状态均衡服务器池，轮询方式调度服务器提供服务
type BalancePool struct {
	servers    map[int32]map[string]*agent.Target
	mu         sync.Mutex
	rand       *rand.Rand
	clientPool agent.TargetPool
}

// GetTarget 如果客户端没有相应类型服务器,根据算法得到类型为target的服务器
func (bp *BalancePool) GetTarget(target int32, client *agent.Target) *agent.Target {
	if ss, ok := bp.servers[target]; ok && len(ss) > 0 {
		keys := reflect.ValueOf(ss).MapKeys()
		return ss[keys[bp.algorithm(len(ss))].String()]
	}
	return nil
}

// SetTargetPool 目标池实现
func (bp *BalancePool) SetTargetPool(tp agent.TargetPool) {
	bp.clientPool = tp
}

var last = 0

//轮询模式选择,TODO有bug
func (bp *BalancePool) algorithm(all int) int {
	if all == 0 {
		return 0
	}
	last = (last + 1) % all
	return last
}

func (bp *BalancePool) hasServer(config *config.ServerConfig) bool {
	if servers, ok := bp.servers[config.Type]; ok {
		_, ok := servers[config.Addr]
		return ok
	}
	return false
}

func (bp *BalancePool) registerServer(cs *rpc.ClientService) {
	go func() {
		err := cs.Init()
		if err == nil {
			if _, ok := bp.servers[cs.Typeid]; !ok {
				bp.servers[cs.Typeid] = make(map[string]*agent.Target)
			}
			server := agent.NewTarget(cs.Stream(), cs.TypeID())
			bp.servers[cs.Typeid][cs.Addr] = server
			logger.Infof("%s(%d) is running", cs.Addr, cs.Typeid)
			err = server.Run(bp.clientPool)
			delete(bp.servers[cs.Typeid], cs.Addr)
		}
		logger.Error(err)
	}()
}

// NewBalancePool 新服务器均衡池
func NewBalancePool() (bp *BalancePool) {
	bp = &BalancePool{}
	bp.rand = rand.New(rand.NewSource(99))
	bp.servers = make(map[int32]map[string]*agent.Target)
	return
}

//从配置频道服务器得到服务器信息
func (bp *BalancePool) subscribe() error {
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
			if !bp.hasServer(info) {
				logger.Info("register ", info)
				bp.registerServer(&rpc.ClientService{Typeid: info.Type, Addr: info.Addr})
			}
		}
	}()
	return nil
}

// Start 启动均衡池，就是订阅服务器信息
func (bp *BalancePool) Start() {
	err := bp.subscribe()
	if err != nil {
		logger.Error(err)
	}
}
