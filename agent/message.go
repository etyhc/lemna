package agent

import (
	"lemna/agent/proto"
	"lemna/arpc"
	"lemna/msg"
)

//MsgBuilder 代理服务器消息构建接口,需要根据不同协议实现
type MsgBuilder interface {
	InvalidTarget() interface{}
}

//MsgKits 代理消息工具集，用于代理服务器构建自己的转发消息
type MsgKits struct {
	//Helper 消息助理
	Helper msg.Helper
	//Builder 消息构建器
	Builder MsgBuilder
}

var msgkits = MsgKits{Helper: msg.ProtoHelper{}, Builder: proto.Builder{}}

//MsgWrapper 将消息封装成转发消息
func MsgWrapper(msg interface{}, target uint32) *arpc.ForwardMsg {
	raw, err := msgkits.Helper.ToRaw(msg)
	if err != nil {
		return nil
	}
	return &arpc.ForwardMsg{Target: target,
		Msg: &arpc.RawMsg{
			Mid: msgkits.Helper.BaseInfo(msg).ID(),
			Raw: raw}}
}
