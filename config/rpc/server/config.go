package server

import "encoding/json"

// Config 服务器配置
type Config struct {
	Name string `json:"name"` //服务器名字
	Addr string `json:"addr"` //服务器地址
	Type int32  `json:"type"` //服务器类型
}

func (c *Config) String() string {
	ret, _ := json.Marshal(*c)
	return string(ret)
}

// FromString 从info中初始化c
func (c *Config) FromString(info string) error {
	return json.Unmarshal([]byte(info), c)
}
