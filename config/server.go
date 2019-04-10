package config

import "encoding/json"

// ServerConfig 服务器配置,用于代理发现服务器
//              服务器启动后发布自己的信息给代理
type ServerConfig struct {
	Name string `json:"name"` //服务器名字
	Addr string `json:"addr"` //服务器地址
	Type int32  `json:"type"` //服务器类型
}

func (sc *ServerConfig) String() string {
	ret, _ := json.Marshal(*sc)
	return string(ret)
}

// FromString 从字串中初始化服务器配置,使用json编码
func (sc *ServerConfig) FromString(info string) error {
	return json.Unmarshal([]byte(info), sc)
}
