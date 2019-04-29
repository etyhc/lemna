package rpc

import (
	"lemna/content"
	"lemna/logger"
	"reflect"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Channel 频道用户
type Channel struct {
	Addr string //频道地址
}

// Publish content.Channel的rpc实现
func (c *Channel) Publish(ctt content.Content) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err == nil {
		sc := NewChannelClient(conn)
		_, err = sc.Publish(context.Background(), &ContentMsg{Info: ctt.ToString(), Name: ctt.Topic()})
		logger.Debug(err, ctt.ToString())
		conn.Close()
	}
	return err
}

// Subscribe content.Channel的rpc实现
func (c *Channel) Subscribe(ctt content.Content) (<-chan content.Content, error) {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	sc := NewChannelClient(conn)
	stream, err := sc.Subscribe(context.Background(), &ContentMsg{Name: ctt.Topic()})
	if err != nil {
		conn.Close()
		return nil, err
	}

	ret := make(chan content.Content)
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				conn.Close()
				break
			}
			content := reflect.New(reflect.TypeOf(ctt).Elem()).Interface().(content.Content)
			err = content.FromString(in.Info)
			if err == nil {
				logger.Debug(content)
				ret <- content
			} else {
				logger.Error(err)
			}
		}
	}()
	return ret, nil
}
