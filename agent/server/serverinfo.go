package server

//服务器调度策略
const (
	SERVERSCHENIL   int32 = iota //不接受调度
	SERVERSCHEROUND              //轮流调度
	SERVERSCHELOAD               //最小负载调度
)

// Info 服务器信息,用于代理发现服务器
//           服务器启动后发布自己的信息给代理
type Info struct {
	Addr string `json:"addr"` //服务器地址
	Type int32  `json:"type"` //服务器类型
	Sche int32  `json:"sche"` //服务器调度策略
	Load int32  `json:"load"` //服务器负载
}

// Topic 服务器信息主题
//       用于服务器发布和代理服务器订阅
func (si *Info) Topic() string {
	return "ServerInfo"
}
