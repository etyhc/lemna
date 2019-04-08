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
	Port     string        //代理地址
	Cm       ClientManager //客户端管理器
	Balancer Balancer      //均衡器
	Token    Token         //Token
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
	sessionid, err := as.Token.GetSessionID(token[0])
	if err != nil {
		return nil, err
	}
	//根据sessionid从client管理器创建一个Client
	_, err = as.Cm.NewClient(sessionid)
	if err != nil {
		return nil, err
	}
	//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
	grpc.SetHeader(ctx, metadata.Pairs("session", fmt.Sprint(sessionid)))
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
	c, err := as.Cm.GetClient(sessionid)
	if err != nil {
		return err
	}
	c.SetStream(stream)
	err = as.runClient(c)
	as.Cm.DelClient(sessionid)
	return err
}

// runClient 接收客户端消息并转发给均衡器提供的服务器
func (as *Service) runClient(c Client) error {
	for {
		cfmsg, err := c.Stream().Recv()
		if err != nil {
			return c.Error(err)
		}
		//转发指令
		server, ok := as.Balancer.GetServer(cfmsg.Target, c)
		if !ok {
			return fmt.Errorf("<target=%d> not find server", cfmsg.Target)
		}

		cfmsg.Target = c.SessionID()
		err = server.Stream().Send(cfmsg)
		if err != nil {
			return server.Error(err)
		}
	}
}

// RunServer 接收服务器消息并转发给相应的客户端
func (as *Service) RunServer(s Server) error {
	for {
		sfmsg, err := s.Stream().Recv()
		if err != nil {
			return s.Error(err)
		}
		client, err := as.Cm.GetClient(sfmsg.Target)
		if err != nil {
			logger.Error(err)
			continue //未找到用户服务器继续服务
		}
		sfmsg.Target = s.TypeID()
		err = client.Stream().Send(sfmsg)
		if err != nil {
			//用户失效，移除用户
			logger.Error(client.Error(err))
			as.Cm.DelClient(client.SessionID())
		}
	}
}

// Run 运行代理服务
func (as *Service) Run() error {
	lis, err := net.Listen("tcp", as.Port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	rpc.RegisterClientServer(s, as)
	return s.Serve(lis)
}