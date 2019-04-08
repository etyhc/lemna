package rpc

import (
	"lemna/config"
	"lemna/logger"
	"reflect"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type ChannelClient struct {
	Addr string
}

func (f *ChannelClient) Publish(s config.Stringer) error {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure())
	if err == nil {
		sc := NewConfigClient(conn)
		_, err = sc.Publish(context.Background(), &ConfigMsg{Info: s.String()})
		logger.Debug(err, s.String())
		conn.Close()
	}
	return err
}

func (f *ChannelClient) Subscribe(info string, t config.Stringer) (<-chan config.Stringer, error) {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	sc := NewConfigClient(conn)
	stream, err := sc.Subscribe(context.Background(), &ConfigMsg{Info: info})
	if err != nil {
		conn.Close()
		return nil, err
	}

	ret := make(chan config.Stringer)
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				conn.Close()
				break
			}
			info := reflect.New(reflect.TypeOf(t).Elem()).Interface().(config.Stringer)
			info.FromString(in.Info)
			logger.Debug(info)
			ret <- info
		}
	}()
	return ret, nil
}