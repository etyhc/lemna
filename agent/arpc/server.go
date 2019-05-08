package arpc

import (
	fmt "fmt"
	"lemna/logger"

	proto "github.com/golang/protobuf/proto"
)

// Server 服务器rpc服务端
type Server struct {
	typeid uint32               //服务器类型
	stream Server_ForwardServer //代理服务器流
	mc     *MsgCenter           //消息中心
	id     uint32               //代理服务器唯一ID
}

// NewServer 接收到代理服务器的rpc调用，生成新的rpc服务
//    typeid 服务器类型，所有服务器类型都一样
//    straem 代理消息流
//        mc 消息处理中心
//        id 代理服务器唯一ID
func NewServer(typeid uint32, stream Server_ForwardServer, mc *MsgCenter, id uint32) *Server {
	return &Server{typeid: typeid, stream: stream, mc: mc, id: id}
}

// Broadcast 广播消息
//   targets 将消息广播给切片里的所有客户端
//       msg 被广播消息
func (s *Server) Broadcast(targets []uint32, msg interface{}) error {
	bmsg, err := s.mc.WrapBM(targets, msg.(proto.Message))
	if err != nil {
		return err
	}
	return s.stream.Send(bmsg)
}

// Forward 转发消息
//         内部实现仍是广播消息，只是广播只有1个客户端
//  target 将消息转发给此客户端
//     msg 被转发消息
func (s *Server) Forward(target uint32, msg interface{}) error {
	bmsg, err := s.mc.WrapBM([]uint32{target}, msg.(proto.Message))
	if err != nil {
		return err
	}
	return s.stream.Send(bmsg)
}

// ID 代理服务器唯一ID
func (s *Server) ID() uint32 {
	return s.id
}

// TypeID 服务器类型ID
func (s *Server) TypeID() uint32 {
	return s.typeid
}

// Error 在错误信息中添加代理服务器信息
func (s *Server) Error(err interface{}) error {
	return fmt.Errorf("<%d> %s", s.id, err)
}

// Handle 接收代理服务器消息并处理消息
func (s *Server) Handle() error {
	for {
		in, err := s.stream.Recv()
		if err == nil {
			err = s.mc.Handle(in, s)
			if err != nil {
				//忽略错误的消息
				logger.Error(err)
			}
		} else {
			return err
		}
	}
}
