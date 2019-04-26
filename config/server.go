package config

import "encoding/json"

// 调度策略
const (
	SERVERSCHENIL   int32 = iota //不接受调度
	SERVERSCHEROUND              //轮流调度
	SERVERSCHELOAD               //最小负载调度
)

//ServerInfo 服务器信息,用于代理发现服务器
//           服务器启动后发布自己的信息给代理
type ServerInfo struct {
	Addr string `json:"addr"` //服务器地址
	Type int32  `json:"type"` //服务器类型
	Sche int32  `json:"sche"` //服务器调度策略
	Load int32  `json:"load"` //服务器负载
}

func (si *ServerInfo) String() string {
	ret, _ := json.Marshal(*si)
	return string(ret)
}

//FromString 从字串中初始化服务器配置,使用json编码
func (si *ServerInfo) FromString(info string) error {
	return json.Unmarshal([]byte(info), si)
}
