package srpc

import (
	context "context"
	"time"

	grpc "google.golang.org/grpc"
)

//Client 服务器调用客户端
type Client struct {
	Addr string //代理服务器地址

	conn   *grpc.ClientConn
	client StoSClient
	ctx    context.Context
}

//Call Call实现
func (c *Client) Call(in *CallMsg) (*CallMsg, error) {
	return c.client.Call(c.ctx, in)
}

//Connect 连接服务器
func (c *Client) Connect() error {
	//连接
	var err error
	c.conn, err = grpc.Dial(c.Addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second*3))
	if err != nil {
		return err
	}
	//rpc客户端
	c.client = NewStoSClient(c.conn)
	c.ctx = context.Background()
	return nil
}

//Close 客户端关闭
func (c *Client) Close() {
	c.conn.Close()
}
