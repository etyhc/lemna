//Package client 基于arpc的封装实现
package client

import (
	"context"
	"lemna/arpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Crpc 客户端rpc基本封装
type Crpc struct {
	Token string //登录代理服务器所需的Token
	Addr  string //代理服务器地址

	stream  arpc.Crpc_ForwardClient
	ostream arpc.Crpc_OtherClient
	conn    *grpc.ClientConn
	client  arpc.CrpcClient
	ctx     context.Context
	header  metadata.MD
}

//GetRequestMetadata 实现credentials.PerRPCCredentials接口
//                   将crpc.Token作为metadata数据放入头中用于代理服务器登录验证
func (crpc *Crpc) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"token": crpc.Token}, nil
}

//RequireTransportSecurity 实现credentials.PerRPCCredentials接口
//                         目前未实现加密传输
func (crpc *Crpc) RequireTransportSecurity() bool {
	return false
}

//Send 发送转发消息给服务器
func (crpc *Crpc) Send(msg *arpc.ForwardMsg) error {
	return crpc.stream.Send(msg)
}

//Recv 接收服务器转发消息
func (crpc *Crpc) Recv() (*arpc.ForwardMsg, error) {
	return crpc.stream.Recv()
}

//AgentSend 发送消息给代理
func (crpc *Crpc) AgentSend(msg *arpc.RawMsg) error {
	return crpc.ostream.Send(msg)
}

//AgentRecv 接收代理消息
func (crpc *Crpc) AgentRecv() (*arpc.RawMsg, error) {
	return crpc.ostream.Recv()
}

//Close 关闭rpc
func (crpc *Crpc) Close() {
	crpc.conn.Close()
}

//Login 登录代理服务器，并发起转发调用
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
		if err == nil {
			crpc.ostream, err = crpc.client.Other(metadata.NewOutgoingContext(crpc.ctx, crpc.header))
		}
	}
	if err != nil {
		crpc.conn.Close()
	}
	return err
}
