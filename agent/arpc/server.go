package arpc

import (
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"
)

type Server struct {
	typeid int32
	stream Server_ForwardServer
	mc     *MsgCenter //消息中心
	id     uint32
}

func NewServer(typeid int32, stream Server_ForwardServer, mc *MsgCenter, id uint32) *Server {
	return &Server{typeid: typeid, stream: stream, mc: mc, id: id}
}

func (s *Server) Broadcast(targets []int32, msg interface{}) error {
	bmsg, err := s.mc.WrapBM(targets, msg.(proto.Message))
	if err != nil {
		return err
	}
	return s.stream.Send(bmsg)
}

func (s *Server) Forward(target int32, msg interface{}) error {
	bmsg, err := s.mc.WrapBM([]int32{target}, msg.(proto.Message))
	if err != nil {
		return err
	}
	return s.stream.Send(bmsg)
}

func (s *Server) ID() uint32 {
	return s.id
}

func (s *Server) TypeID() int32 {
	return s.typeid
}

func (s *Server) Error(err interface{}) error {
	return fmt.Errorf("<%d> %s", s.typeid, err)
}
