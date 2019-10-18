//Package server 服务器对代理服务器发起的rpc调用封装
package server

import (
	"context"
	"encoding/json"
	"lemna/agent/server"
	"lemna/arpc"
	"lemna/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//Srpc 服务器rpc封装
type Srpc struct {
	Addr    string //代理服务器地址
	Info    server.Info
	fstream arpc.Srpc_ForwardClient
	mstream arpc.Srpc_MulticastClient
	ostream arpc.Srpc_OtherClient
	conn    *grpc.ClientConn
}

//Send 发送转发消息给客户端
func (c *Srpc) Send(msg *arpc.ForwardMsg) error {
	return c.fstream.Send(msg)
}

//Recv 接收客户端转发消息
func (c *Srpc) Recv() (*arpc.ForwardMsg, error) {
	return c.fstream.Recv()
}

//Multicast 向多个客户端发送转发消息
func (c *Srpc) Multicast(msg *arpc.MulticastMsg) error {
	return c.mstream.Send(msg)
}

//AgentSend 向代理服务器发送消息
func (c *Srpc) AgentSend(msg *arpc.RawMsg) error {
	return c.ostream.Send(msg)
}

//AgentRecv 接收代理服务器消息
func (c *Srpc) AgentRecv() (*arpc.RawMsg, error) {
	return c.ostream.Recv()
}

//Close 关闭链接
func (c *Srpc) Close() {
	c.conn.Close()
}

//Connect 链接代理服务器
//        并发起Forward、Multicast、Other 远程调用
func (c *Srpc) Connect() error {
	b, err := json.Marshal(c.Info)
	if err != nil {
		return err
	}
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("info", string(b)))

	//连接
	c.conn, err = grpc.Dial(c.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return err
	}
	//rpc客户端
	ac := arpc.NewSrpcClient(c.conn)
	//发起rpc调用
	if err == nil {
		c.fstream, err = ac.Forward(ctx)
		if err == nil {
			c.mstream, err = ac.Multicast(ctx)
			if err == nil {
				c.ostream, err = ac.Other(ctx)
			}
		}
	}
	if err != nil {
		c.conn.Close()
	} else {
		logger.Infof("Connect agent(%s) success.", c.Addr)
	}
	return err
}
