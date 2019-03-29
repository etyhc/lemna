package agent

import (
	"errors"
	fmt "fmt"
	"lemna/agent/rpc"
	"net"
	"strconv"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AgentService ...
type AgentService struct {
	Port     string
	Cm       ClientManager
	Balancer Balancer
	Token    Token
}

func (as *AgentService) Register(cont context.Context, msg *rpc.ClientRegMsg) (ret *rpc.ClientRegMsg, err error) {
	ret = msg
	var sessionid int32
	//根据token分配一个sessionid
	if sessionid, err = as.Token.GetSessionID(msg.Token); err == nil {
		var c Client
		//根据sessionid从client管理器得到或创建一个Client
		if c, err = as.Cm.GetClient(sessionid); err == nil {
			//将session返回给客户端，客户端每次RPC调用都应将此session放入head中
			grpc.SetHeader(cont, metadata.Pairs("session", fmt.Sprint(c.SessionID())))
		}
	}
	return
}

func (as *AgentService) Forward(stream rpc.Client_ForwardServer) (err error) {

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		if session, ok := md["session"]; ok {
			tmp, _ := strconv.Atoi(session[0])
			sessionid := int32(tmp)
			var c Client
			if c, err = as.Cm.GetClient(sessionid); err == nil {
				c.SetStream(stream)
				err = c.Run()
				as.Cm.DelClient(c.SessionID())
			}
			return
		}
	}
	err = errors.New("invalid client,no session")
	return
}

func (as *AgentService) Start() error {
	lis, err := net.Listen("tcp", as.Port)
	if err == nil {
		s := grpc.NewServer()
		rpc.RegisterClientServer(s, as)
		err = s.Serve(lis)
	}
	return err
}
