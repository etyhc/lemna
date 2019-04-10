package rpc

import (
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ConfigServerAddr 频道服务器默认地址
var ConfigServerAddr = ":10000"

// ChannelService 一个基于grpc的配置订阅/发布频道服务
type ChannelService struct {
	subscribers map[string]Channel_SubscribeServer
	topics      map[string]map[string]int
	addr        string
}

// NewChannelService 新的配置频道服务
func NewChannelService(addr string) *ChannelService {
	return &ChannelService{
		subscribers: make(map[string]Channel_SubscribeServer),
		topics:      make(map[string]map[string]int),
		addr:        addr}
}

// Publish rpc配置频道发布实现
func (ch *ChannelService) Publish(ctx context.Context, msg *ConfigMsg) (*ConfigMsg, error) {
	//新主题加入
	topic, ok := ch.topics[msg.Name]
	if !ok {
		topic = make(map[string]int)
		ch.topics[msg.Name] = topic
	}
	topic[msg.Info] = topic[msg.Info] + 1
	logger.Debug("pub<", msg.Name, ":", msg.Info, ">")
	//新主题发送给订阅者
	for addr := range ch.subscribers {
		err := ch.subscribers[addr].Send(msg)
		if err == nil {
			logger.Debug("sub: ...")
		} else {
			logger.Error(err)
			delete(ch.subscribers, addr)
		}
	}
	return msg, nil
}

// Subscribe rpc配置频道订阅实现
func (ch *ChannelService) Subscribe(msg *ConfigMsg, stream Channel_SubscribeServer) error {
	//发送订阅主题给订阅者
	for info := range ch.topics[msg.Name] {
		msg.Info = info
		err := stream.Send(msg)
		if err != nil {
			return err
		}
		logger.Debug("sub: ", msg.Info)
	}

	//订阅者加入
	if p, ok := peer.FromContext(stream.Context()); ok {
		ch.subscribers[p.Addr.String()] = stream
	}
	//等待订阅者失效
	<-stream.Context().Done()
	return nil
}

// Run 运行配置频道服务
func (ch *ChannelService) Run() error {
	lis, err := net.Listen("tcp", ch.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	RegisterChannelServer(s, ch)
	return s.Serve(lis)
}
