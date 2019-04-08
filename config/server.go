package config

import "encoding/json"

// ServerConfig 服务器配置
type ServerConfig struct {
	Name string `json:"name"` //服务器名字
	Addr string `json:"addr"` //服务器地址
	Type int32  `json:"type"` //服务器类型
}

func (sc *ServerConfig) String() string {
	ret, _ := json.Marshal(*sc)
	return string(ret)
}

// FromString 从info中初始化c
func (sc *ServerConfig) FromString(info string) error {
	return json.Unmarshal([]byte(info), sc)
}
