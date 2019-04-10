package agent

import (
	fmt "fmt"
	"lemna/agent/rpc"
	"lemna/logger"
	"net"
	"strconv"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Service 代理服务
type Service struct {
	addr       string //代理地址
	serverPool TargetPool
	token      Token //Token
	clientPool *clientManager
}

func NewService(addr string, serverPool TargetPool, t Token) *Service {
	cp := newClientMananger()
	cp.SetTargetPool(serverPool)
	serverPool.SetTargetPool(cp)
	return &Service{addr: addr, serverPool: serverPool, token: t, clientPool: cp}
}

// Register rpc.ClientServer.Register的实现,根据token分配唯一sessionid，并将此ID通过消息头返回给客户端
func (as *Service) Register(ctx context.Context, msg *rpc.ClientRegMsg) (*rpc.ClientRegMsg, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("invalid rpc,no metadata")
	}
	token, ok := md["token"]
	if !ok {
		return nil, fmt.Errorf("invalid client,no token")
	}
	//通过验证从ctx中获得sessionid
	sessionid, err := as.token.GetSessionID(token[0])
	if err != nil {
		return nil, err
	}
	//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
	err = grpc.SetHeader(ctx, metadata.Pairs("session", fmt.Sprint(sessionid)))
	if err != nil {
		return nil, err
	}
	logger.Debug(sessionid, " Register")
	return msg, nil
}

// Forward rpc.ClientServer.Forward的实现,会验证消息头的session数据是否有效
func (as *Service) Forward(stream rpc.Client_ForwardServer) error {
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
	sessionid := int32(tmp)
	//根据sessionid从client管理器创建一个Client
	client := as.clientPool.newClient(stream, sessionid)
	err := client.Run(as.clientPool.serverPool)
	as.clientPool.delClient(sessionid)
	return err
}

// Run 运行代理服务
func (as *Service) Run() error {
	lis, err := net.Listen("tcp", as.addr)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	rpc.RegisterClientServer(s, as)
	return s.Serve(lis)
}
