package rpc

import (
	"lemna/config"
	"lemna/logger"
	"reflect"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type ChannelUser struct {
	Addr string
}

func (f *ChannelUser) Publish(s config.Stringer) error {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure())
	if err == nil {
		sc := NewChannelClient(conn)
		_, err = sc.Publish(context.Background(), &ConfigMsg{Info: s.String()})
		logger.Debug(err, s.String())
		conn.Close()
	}
	return err
}

func (f *ChannelUser) Subscribe(info string, t config.Stringer) (<-chan config.Stringer, error) {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	sc := NewChannelClient(conn)
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
			err = info.FromString(in.Info)
			if err == nil {
				logger.Debug(info)
				ret <- info
			} else {
				logger.Error(err)
			}
		}
	}()
	return ret, nil
}
