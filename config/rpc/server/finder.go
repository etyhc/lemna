package server

import (
	"lemna/config"
	"lemna/config/rpc"
	"lemna/logger"
	"reflect"

	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

type Finder struct {
	Addr string
}

func (f *Finder) Publish(s config.Serialization) error {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure())
	if err == nil {
		sc := rpc.NewConfigClient(conn)
		_, err = sc.Publish(context.Background(), &rpc.ConfigMsg{Info: s.String()})
		logger.Debug(err, s.String())
		conn.Close()
	}
	return err
}

func (f *Finder) Subscribe(info string, t config.Serialization) (<-chan config.Serialization, error) {
	conn, err := grpc.Dial(f.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	sc := rpc.NewConfigClient(conn)
	stream, err := sc.Subscribe(context.Background(), &rpc.ConfigMsg{Info: info})
	if err != nil {
		conn.Close()
		return nil, err
	}

	ret := make(chan config.Serialization)
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil {
				conn.Close()
				break
			}
			info := reflect.New(reflect.TypeOf(t).Elem()).Interface().(config.Serialization)
			info.Init(in.Info)
			logger.Debug(info)
			ret <- info
		}
	}()
	return ret, nil
}
