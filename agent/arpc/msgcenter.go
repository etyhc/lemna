package arpc

import (
	"fmt"
	"lemna/utils"
	"reflect"

	proto "github.com/golang/protobuf/proto"
)

// MsgStream 消息流，用于接收和发送消息
//           MsgCenter处理消息时会将消息流代入消息回调函数，以便消息回复.
//           这个接口只有一个用处，让服务器程序写游戏客户端测试代码能够复用msgcenter代码
//           否则可以使用*Server
type MsgStream interface {
	// Broadcast 服务器向客户端广播消息
	Broadcast([]uint32, interface{}) error
	// Forward 向服务器或者客户端转发消息
	Forward(uint32, interface{}) error
	// ID 消息流唯一ID
	ID() uint32
}

// MsgHandler 是个消息回调函数，需要实现，并注册到MsgCenter
//      int32 消息来源
//  interface 消息本体
//  MsgServer 消息来源流，用于回复
type MsgHandler func(uint32, interface{}, MsgStream)

type msgInfo struct {
	elem    reflect.Type
	handler MsgHandler
}

/*MsgCenter 消息中心,封装消息工具*/
type MsgCenter struct {
	hash map[string]uint32
	info map[uint32]msgInfo
}

/*NewMsgCenter 新消息中心*/
func NewMsgCenter() *MsgCenter {
	return &MsgCenter{hash: make(map[string]uint32), info: make(map[uint32]msgInfo)}
}

// Reg 消息注册
//      msg 被注册的消息
//  handler 消息处理函数
func (mc *MsgCenter) Reg(msg interface{}, handler MsgHandler) {
	name := reflect.TypeOf(msg).Elem().Name()
	hash := utils.HashFnv1a(name)
	if info, ok := mc.info[hash]; ok {
		panic(fmt.Sprintf("Hash(%d) conflict %s %s", hash, name, info.elem.Name()))
	}
	mc.info[hash] = msgInfo{reflect.TypeOf(msg).Elem(), handler}
	mc.hash[name] = hash
}

// WrapFMNoCheck 将消息封装成转发消息,无需提前注册此消息
//        target 转发目标
//           msg 被封装的消息
func WrapFMNoCheck(target uint32, msg proto.Message) (*ForwardMsg, error) {
	hash := utils.HashFnv1a(reflect.TypeOf(msg).Elem().Name())
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &ForwardMsg{Target: target, Msg: &RawMsg{Type: hash, Raw: buf}}, nil
}

func (mc *MsgCenter) wrapRaw(msg proto.Message) (*RawMsg, error) {
	hash, ok := mc.hash[reflect.TypeOf(msg).Elem().Name()]
	if !ok {
		return nil, fmt.Errorf("%s don't register", reflect.TypeOf(msg).Elem().Name())
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return &RawMsg{Type: hash, Raw: buf}, nil
}

// WrapFM 将消息封装为转发消息,此消息需提前注册
//  target 转发目标
//     msg 被转发的消息
func (mc *MsgCenter) WrapFM(target uint32, msg proto.Message) (*ForwardMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &ForwardMsg{Target: target, Msg: raw}, nil
	}
	return nil, err
}

// WrapBM 将消息封装为广播消息,此消息需提前注册
//  targets 转发目标切片
//      msg 被广播的消息
func (mc *MsgCenter) WrapBM(targets []uint32, msg proto.Message) (*BroadcastMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &BroadcastMsg{Targets: targets, Msg: raw}, nil
	}
	return nil, err
}

// Handle 转发消息处理函数
//        转发消息会被解码成原始消息，并调用注册过的处理函数来处理此消息
//   fmsg 收到的转发消息
//   from 消息流
func (mc *MsgCenter) Handle(fmsg *ForwardMsg, from MsgStream) error {
	info, ok := mc.info[fmsg.Msg.Type]
	if !ok {
		return fmt.Errorf("unregistered message type %d", fmsg.Msg.Type)
	}
	msg := reflect.New(info.elem).Interface()
	err := proto.Unmarshal(fmsg.Msg.Raw, msg.(proto.Message))
	if err != nil {
		return err
	}
	info.handler(fmsg.Target, msg, from)
	return nil
}
