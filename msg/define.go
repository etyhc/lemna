package msg

//RawMsg 普通消息接口
type RawMsg interface {
	GetMid() uint32
	GetRaw() []byte
}

//ForwardMsg 转发消息接口
type ForwardMsg interface {
	GetTarget() uint32
	RawMsg
}

//MulticastMsg 广播消息接口
type MulticastMsg interface {
	GetTargets() []uint32
	RawMsg
}
