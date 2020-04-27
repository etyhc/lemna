package srpc

//MsgInfo 消息基本信息
type MsgInfo interface {
	ID() uint32   //消息ID，用于区分不同消息，客户端和服务器协商确定
	Name() string //消息名字
}

//Helper 消息辅助接口，不同消息协议自行实现
type Helper interface {
	//Extract 提取消息基本信息
	Extract(interface{}) MsgInfo
	//ToRaw 序列化消息
	ToRaw(*CallMsg) (uint32, []byte, error)
	//FromRaw 反序列化消息
	FromRaw(MsgInfo, []byte) (*CallMsg, error)
}

// Handler 是个消息回调函数，需要实现，并注册到MsgCenter
//      int32 消息来源
//  interface 消息本体
//     Stream 消息来源流，用于回复
type Handler func(*CallMsg) (*CallMsg, error)
