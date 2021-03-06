// Package client 实现了客户端的rpc服务.
package client

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/arpc"
	"lemna/logger"
	"net"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Service 代理服务，接受客户端连接，并验证
//         将客户端消息转发给服务器并将服务器消息转发给客户端
type Service struct {
	addr      string           //代理地址
	token     Token            //Token
	clientmgr *manager         //客户端管理器
	sp        agent.TargetPool //服务器池
}

// NewService 新代理服务
func NewService(addr string, t Token) *Service {
	return &Service{addr: addr, token: t, clientmgr: newMananger()}
}

// GetTarget 目标池接口实现
func (cs *Service) GetTarget(cid uint32) agent.Target {
	return cs.clientmgr.getTarget(cid)
}

// Login rpc.ClientServer.Login接口实现
//       根据token分配唯一sessionid，并将此ID通过消息头返回给客户端
//       客户端调用Forward时应将此头返回给服务器
func (cs *Service) Login(ctx context.Context, msg *arpc.LoginMsg) (*arpc.LoginMsg, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid rpc,no metadata")
	}
	token, ok := md["token"]
	if !ok {
		return nil, fmt.Errorf("invalid client,no token")
	}
	//通过验证从ctx中获得sessionid
	session, err := cs.token.GetSession(token[0])
	if err != nil {
		return nil, err
	}
	//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
	err = grpc.SetHeader(ctx, metadata.Pairs("session", session))
	if err != nil {
		return nil, err
	}
	return msg, nil

}

func getUID(ctx context.Context, token Token) (uint32, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("invalid rpc,no metadata")
	}
	session, ok := md["session"]
	if !ok {
		return 0, fmt.Errorf("invalid client,no session")
	}
	return token.GetUID(session[0])
}

// Call rpc.ClientServer.Call接口实现
func (cs *Service) Call(context.Context, *arpc.RawMsg) (*arpc.RawMsg, error) {
	//TODO 还未实现,处理来自客户端的Call
	return nil, nil
}

// Forward rpc.ClientServer.Forward接口实现
//         会验证消息头的session数据是否有效
func (cs *Service) Forward(stream arpc.CAgent_ForwardServer) error {
	uid, err := getUID(stream.Context(), cs.token)
	if err != nil {
		return err
	}
	logger.Infof("<uid=%d> Logining", uid)
	//根据uid从client管理器初始化一个Client
	client, err := cs.clientmgr.newTarget(stream, uid)
	if err == nil {
		err = client.Forward(cs.sp)
		cs.clientmgr.delTarget(uid)
	}
	logger.Errorf("%s", err)
	return err
}

// Run 运行代理服务,接受客户端的连接
func (cs *Service) Run(sp agent.TargetPool) error {
	cs.sp = sp
	lis, err := net.Listen("tcp", cs.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	arpc.RegisterCAgentServer(s, cs)
	logger.Infof("Start client service at %s", cs.addr)
	return s.Serve(lis)
}
