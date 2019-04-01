package agent

import (
	"errors"
	fmt "fmt"
	"lemna/rpc"
	"net"
	"strconv"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AgentService 代理服务
type AgentService struct {
	Port     string
	Cm       ClientManager
	Balancer Balancer
	Token    Token
}

// Register rpc.ClientServer.Register的实现,根据token分配唯一sessionid，并将此ID通过消息头返回给客户端
func (as *AgentService) Register(cont context.Context, msg *rpc.ClientRegMsg) (ret *rpc.ClientRegMsg, err error) {
	ret = msg
	var sessionid int32
	//根据token分配一个sessionid
	if sessionid, err = as.Token.GetSessionID(msg.Token); err == nil {
		//根据sessionid从client管理器创建一个Client
		if _, err = as.Cm.NewClient(sessionid); err == nil {
			//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
			grpc.SetHeader(cont, metadata.Pairs("session", fmt.Sprint(sessionid)))
		}
	}
	return
}

// Forward rpc.ClientServer.Forward的实现,会验证消息头的session数据是否有效
func (as *AgentService) Forward(stream rpc.Client_ForwardServer) (err error) {

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		if session, ok := md["session"]; ok {
			tmp, _ := strconv.Atoi(session[0])
			sessionid := int32(tmp)
			var c Client
			if c, err = as.Cm.GetClient(sessionid); err == nil {
				c.SetStream(stream)
				err = as.runClient(c)
				as.Cm.DelClient(sessionid)
			}
			return
		}
	}
	err = errors.New("invalid client,no session")
	return
}

// runClient 接收客户端消息并转发给均衡器提供的服务器
func (as *AgentService) runClient(c Client) (err error) {
	var cfmsg *rpc.ForwardMsg
	for {
		cfmsg, err = c.Stream().Recv()
		if err == nil {
			//转发指令
			if server, ok := as.Balancer.GetServer(cfmsg.Target, c.SessionID()); ok {
				cfmsg.Target = c.SessionID()
				err = server.Stream().Send(cfmsg)
			} else {
				err = fmt.Errorf("not find server. target=%d", cfmsg.Target)
			}
		}
		if err != nil {
			return
		}
	}
}

// RunServer 接收服务器消息并转发给相应的客户端
func (as *AgentService) RunServer(s Server) (err error) {
	var sfmsg *rpc.ForwardMsg
	for {
		sfmsg, err = s.Stream().Recv()
		if err == nil {
			if client, err := as.Cm.GetClient(sfmsg.Target); err == nil {
				sfmsg.Target = s.TypeID()
				client.Stream().Send(sfmsg)
			}
		} else {
			err = fmt.Errorf("not find client. session=%d", sfmsg.Target)
		}
		if err != nil {
			return
		}
	}
}

// Start 启动代理服务
func (as *AgentService) Start() error {
	lis, err := net.Listen("tcp", as.Port)
	if err == nil {
		s := grpc.NewServer()
		rpc.RegisterClientServer(s, as)
		err = s.Serve(lis)
	}
	return err
}
