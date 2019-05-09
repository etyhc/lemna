// Package crpc 基于grpc的内容发布/订阅实现.
//
//              发布内容是持久的，可以订阅到订阅之前发布的内容
package crpc

import (
	"lemna/content"
	"lemna/logger"
	"reflect"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Channel 基于rpc的内容频道实现
type Channel struct {
	Addr string //频道地址
}

// Publish content.Channel的rpc实现
func (c *Channel) Publish(ctt content.Content) error {
	conn, err := grpc.Dial(c.Addr, grpc.WithInsecure())
	if err == nil {
		defer conn.Close()
		sc := NewChannelClient(conn)
		var ctb []byte
		ctb, err = content.ToJSON(ctt)
		if err == nil {
			_, err = sc.Publish(context.Background(), &ContentMsg{Info: string(ctb), Name: ctt.Topic()})
		}
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
			c := reflect.New(reflect.TypeOf(ctt).Elem()).Interface().(content.Content)
			c, err = content.FromJSON(c, []byte(in.Info))
			logger.Debug(in.Info)
			if err == nil {
				ret <- c
			} else {
				logger.Error(err)
			}
		}
	}()
	return ret, nil
}
