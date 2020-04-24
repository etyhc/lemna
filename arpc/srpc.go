//Package arpc 服务器对代理服务器发起的rpc调用封装
package arpc

import (
	"context"
	"lemna/logger"

	"google.golang.org/grpc"
)

//Srpc 服务器rpc封装
type Srpc struct {
	Addr    string //代理服务器地址
	fstream SAgent_ForwardClient
	mstream SAgent_MulticastClient
	client  SAgentClient
	ctx     context.Context
	conn    *grpc.ClientConn
}

//Send 发送转发消息给客户端
func (c *Srpc) Send(msg *ForwardMsg) error {
	return c.fstream.Send(msg)
}

//Recv 接收客户端转发消息
func (c *Srpc) Recv() (*ForwardMsg, error) {
	return c.fstream.Recv()
}

//Multicast 向多个客户端发送转发消息
func (c *Srpc) Multicast(msg *MulticastMsg) error {
	return c.mstream.Send(msg)
}

//Call 向代理服务器发送消息
func (c *Srpc) Call(msg *RawMsg) (*RawMsg, error) {
	return c.client.Call(c.ctx, msg)
}

//Close 关闭链接
func (c *Srpc) Close() {
	c.conn.Close()
}

//Connect 链接代理服务器
//        并发起Forward、Multicast、Other 远程调用
func (c *Srpc) Connect() error {
	/*
		b, err := json.Marshal(c.Info)
		if err != nil {
			return err
		}
		ctx := metadata.NewOutgoingContext(
			context.Background(),
			metadata.Pairs("info", string(b)))
	*/
	c.ctx = context.Background()
	var err error
	//连接
	c.conn, err = grpc.Dial(c.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return err
	}
	//rpc客户端
	c.client = NewSAgentClient(c.conn)
	//发起rpc调用
	if err == nil {
		c.fstream, err = c.client.Forward(c.ctx)
		if err == nil {
			c.mstream, err = c.client.Multicast(c.ctx)
		}
	}
	if err != nil {
		c.conn.Close()
	} else {
		logger.Infof("Connect agent(%s) success.", c.Addr)
	}
	return err
}
