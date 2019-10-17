package client

import (
	"context"
	"lemna/arpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Crpc struct {
	Token string
	Addr  string

	stream arpc.Crpc_ForwardClient
	conn   *grpc.ClientConn
	client arpc.CrpcClient
	ctx    context.Context
	header metadata.MD
}

func (crpc *Crpc) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"token": crpc.Token}, nil
}

func (crpc *Crpc) RequireTransportSecurity() bool {
	return false
}

func (crpc *Crpc) Send(msg *arpc.ForwardMsg) error {
	return crpc.stream.Send(msg)
}

func (crpc *Crpc) Recv() (*arpc.ForwardMsg, error) {
	return crpc.stream.Recv()
}

func (crpc *Crpc) Close() {
	crpc.conn.Close()
}

func (crpc *Crpc) Login() error {
	//连接
	var err error
	crpc.conn, err = grpc.Dial(crpc.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(crpc))
	if err != nil {
		return err
	}
	//rpc客户端
	crpc.client = arpc.NewCrpcClient(crpc.conn)
	//登录rpc调用
	_, err = crpc.client.Login(crpc.ctx, &arpc.LoginMsg{Token: crpc.Token}, grpc.Header(&crpc.header))
	if err == nil {
		crpc.stream, err = crpc.client.Forward(metadata.NewOutgoingContext(crpc.ctx, crpc.header))
	}
	if err != nil {
		crpc.conn.Close()
	}
	return err
}
