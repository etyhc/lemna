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

// Stream 消息流，用于接收和发送消息
//        MsgProc处理消息时会将消息流代入消息回调函数，以便消息回复.
//        这个接口只有一个用处，让服务器程序写游戏客户端测试代码能够复用msgproc代码
type Stream interface {
	// Forward 向服务器或者客户端转发消息
	Send(uint32, interface{}) error
}

// Handler 是个消息回调函数，需要实现，并注册到MsgCenter
//      int32 消息来源
//  interface 消息本体
//     Stream 消息来源流，用于回复
type Handler func(uint32, interface{}, Stream)

//Info 消息基本信息
type Info interface {
	ID() uint32   //消息ID，用于区分不同消息，客户端和服务器协商确定
	Name() string //消息名字
}

//Helper 消息辅助接口，不同消息协议自行实现
type Helper interface {
	//Extract 提取消息基本信息
	Extract(interface{}) Info
	//ToRaw 序列化消息
	ToRaw(interface{}) (uint32, []byte, error)
	//FromRaw 反序列化消息
	FromRaw(Info, []byte) (interface{}, error)
	Wrap(uint32, interface{}) (ForwardMsg, error)
	WrapMM([]uint32, interface{}) (MulticastMsg, error)
}
