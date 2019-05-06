package arpc

import (
	fmt "fmt"

	proto "github.com/golang/protobuf/proto"
	"google.golang.org/grpc/peer"
)

type Server struct {
	typeid int32
	stream Server_ForwardServer
	mc     *MsgCenter //消息中心
}

func NewServer(typeid int32, stream Server_ForwardServer, mc *MsgCenter) *Server {
	return &Server{typeid: typeid, stream: stream, mc: mc}
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

func (s *Server) GetPeerAddr() (string, bool) {
	peer, ok := peer.FromContext(s.stream.Context())
	if ok {
		return peer.Addr.String(), true
	}
	return "", false
}

func (s *Server) TypeID() int32 {
	return s.typeid
}

func (s *Server) Error(err interface{}) error {
	return fmt.Errorf("<%d> %s", s.typeid, err)
}
