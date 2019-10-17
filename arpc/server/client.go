package client

import (
	"context"
	"lemna/arpc"
	"lemna/logger"

	"google.golang.org/grpc"
)

type Client struct {
	Addr            string
	forwardStream   arpc.Arpc_ForwardClient
	multicastStream arpc.Arpc_MulticastClient
	conn            *grpc.ClientConn
}

func (c *Client) Send(msg *arpc.ForwardMsg) error {
	return c.forwardStream.Send(msg)
}

func (c *Client) Recv() (*arpc.ForwardMsg, error) {
	return c.forwardStream.Recv()
}

func (c *Client) Multicast(msg *arpc.MulticastMsg) error {
	return c.multicastStream.Send(msg)
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Connect() error {
	//连接
	var err error
	c.conn, err = grpc.Dial(c.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		return err
	}
	//rpc客户端
	ac := arpc.NewArpcClient(c.conn)
	ctx := context.Background()
	//发起rpc调用
	if err == nil {
		c.forwardStream, err = ac.Forward(ctx)
		if err == nil {
			c.multicastStream, err = ac.Multicast(ctx)
		}
	}
	if err != nil {
		c.conn.Close()
	} else {
		logger.Info("agent is alive.")
	}
	return err
}
