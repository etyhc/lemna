package crpc

import (
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

// ChannelService 一个基于grpc的内容订阅/发布频道服务
//                内容是持久的，所以订阅者可以接收到订阅之前的发布内容
type ChannelService struct {
	subscribers map[string]Channel_SubscribeServer
	topics      map[string]map[string]int
	addr        string
}

// NewChannelService 新的内容频道服务
func NewChannelService(addr string) *ChannelService {
	return &ChannelService{
		subscribers: make(map[string]Channel_SubscribeServer),
		topics:      make(map[string]map[string]int),
		addr:        addr}
}

// Publish rpc内容频道发布实现
func (ch *ChannelService) Publish(ctx context.Context, msg *ContentMsg) (*ContentMsg, error) {
	//新主题加入
	topic, ok := ch.topics[msg.Name]
	if !ok {
		topic = make(map[string]int)
		ch.topics[msg.Name] = topic
	}
	topic[msg.Info]++
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

// Subscribe rpc内容频道订阅实现
func (ch *ChannelService) Subscribe(msg *ContentMsg, stream Channel_SubscribeServer) error {
	//发送订阅主题给订阅者
	for topic := range ch.topics[msg.Name] {
		msg.Info = topic
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

// Run 运行内容频道服务
func (ch *ChannelService) Run() error {
	lis, err := net.Listen("tcp", ch.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	RegisterChannelServer(s, ch)
	return s.Serve(lis)
}
