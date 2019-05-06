package client

import (
	fmt "fmt"
	"lemna/agent"
	"lemna/agent/arpc"
	"lemna/logger"
	"net"
	"strconv"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// ClientService 代理服务，接受客户端连接，并验证
//         将客户端消息转发给服务器并将服务器消息转发给客户端
type ClientService struct {
	addr      string         //代理地址
	token     Token          //Token
	clientmgr *clientManager //客户端池
	sp        agent.TargetPool
}

// NewClientService 新代理服务
func NewClientService(addr string, t Token) *ClientService {
	cp := newClientMananger()
	return &ClientService{addr: addr, token: t, clientmgr: cp}
}

// GetClient 目标池接口实现
func (cs *ClientService) GetTarget(cid int32) agent.Target {
	ret := cs.clientmgr.getClient(cid)
	if ret != nil {
		return ret
	}
	return nil
}

// SetServerPool 目标池接口实现
func (cs *ClientService) Bind(sp agent.TargetPool) {
	cs.sp = sp
}

// Login rpc.ClientServer.Login接口实现
//       根据token分配唯一sessionid，并将此ID通过消息头返回给客户端
//       客户端调用Forward时应将此头返回给服务器
func (cs *ClientService) Login(ctx context.Context, msg *arpc.LoginMsg) (*arpc.LoginMsg, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid rpc,no metadata")
	}
	token, ok := md["token"]
	if !ok {
		return nil, fmt.Errorf("invalid client,no token")
	}
	//通过验证从ctx中获得sessionid
	sessionid, err := cs.token.GetSessionID(token[0])
	if err != nil {
		return nil, err
	}
	//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
	err = grpc.SetHeader(ctx, metadata.Pairs("session", fmt.Sprint(sessionid)))
	if err != nil {
		return nil, err
	}
	return msg, nil

}

// Forward rpc.ClientServer.Forward接口实现
//         会验证消息头的session数据是否有效
func (cs *ClientService) Forward(stream arpc.Client_ForwardServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return fmt.Errorf("invalid rpc,no metadata")
	}
	session, ok := md["session"]
	if !ok {
		return fmt.Errorf("invalid client,no session")
	}
	logger.Debug(session)
	tmp, _ := strconv.Atoi(session[0])
	uid, err := cs.token.GetUID(int32(tmp))
	if err != nil {
		return err
	}
	//根据sessionid从client管理器初始化一个Client
	client, err := cs.clientmgr.newClient(stream, uid)
	if err == nil {
		err = agent.CtoS(client, cs.sp)
		cs.clientmgr.delClient(uid)
		//通知服务器用户失效
		client.SayBye()
	}
	return err
}

// Run 运行代理服务,接受客户端的连接
func (cs *ClientService) Run() error {
	lis, err := net.Listen("tcp", cs.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	arpc.RegisterClientServer(s, cs)
	return s.Serve(lis)
}