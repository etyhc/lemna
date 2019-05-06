package arpc

import (
	"fmt"
	"hash/fnv"
	"reflect"

	proto "github.com/golang/protobuf/proto"
)

/* MsgServer 消息服务器是接收消息的IO
   MsgCenter处理消息时会将消息服务器传入，以便消息回复
	 这个接口只有一个用处，让服务器程序写游戏客户端测试代码能够使用msgcenter代码
	 否则可以使用*Server*/
type MsgServer interface {
	// Broadcast 服务器向客户端广播消息
	Broadcast([]int32, interface{}) error
	// Forward 向服务器或者客户端转发消息
	Forward(int32, interface{}) error
	ID() uint32
}

/*MsgHandler 是个消息回调函数，需要实现，并注册到MsgCenter
  int32 消息来源
  interface{} 消息本体
	MsgAgent 消息代理*/
type MsgHandler func(int32, interface{}, MsgServer)

type msgInfo struct {
	elem    reflect.Type
	handler MsgHandler
}

/*MsgCenter 消息中心,封装消息工具*/
type MsgCenter struct {
	hash map[string]int32
	info map[int32]msgInfo
}

/*NewMsgCenter 新消息中心*/
func NewMsgCenter() *MsgCenter {
	return &MsgCenter{hash: make(map[string]int32), info: make(map[int32]msgInfo)}
}

func fnvhash(str string) int32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return int32(h.Sum32())
}

/*Reg 消息注册
  msg 消息定义
  handler 消息处理函数*/
func (mc *MsgCenter) Reg(msg interface{}, handler MsgHandler) {
	name := reflect.TypeOf(msg).Elem().Name()
	hash := fnvhash(name)
	if info, ok := mc.info[hash]; ok {
		panic(fmt.Errorf("Hash(%d) conflict %s %s", hash, name, info.elem.Name()))
	}
	mc.info[hash] = msgInfo{reflect.TypeOf(msg).Elem(), handler}
	mc.hash[name] = hash
}

func WrapFMNoCheck(target int32, msg proto.Message) (*ForwardMsg, error) {
	hash := fnvhash(reflect.TypeOf(msg).Elem().Name())
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

/*Wrap 将消息封装为转发消息
  target 转发目标
  msg 被编码的消息*/
func (mc *MsgCenter) WrapFM(target int32, msg proto.Message) (*ForwardMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &ForwardMsg{Target: target, Msg: raw}, nil
	}
	return nil, err
}

/*WrapBroadcast 将消息封装为转发消息
  targets 转发目标数组
  msg 被编码的消息*/
func (mc *MsgCenter) WrapBM(targets []int32, msg proto.Message) (*BroadcastMsg, error) {
	raw, err := mc.wrapRaw(msg)
	if err == nil {
		return &BroadcastMsg{Targets: targets, Msg: raw}, nil
	}
	return nil, err
}

/*Handle 转发消息处理函数
  fmsg 转发消息*/
func (mc *MsgCenter) Handle(fmsg *ForwardMsg, from MsgServer) error {
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