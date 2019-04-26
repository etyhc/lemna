package config

import "encoding/json"

// AgentInfo 代理服务器信息
//           代理服务器启动后，发布自己的信息
type AgentInfo struct {
	AgentID int32  `json:"aid"`
	Load    int32  `json:"load"`
	Addr    string `jsong:"addr"`
}

func (sc *AgentInfo) String() string {
	ret, _ := json.Marshal(*sc)
	return string(ret)
}

// FromString 从字串中初始化服务器配置,使用json编码
func (sc *AgentInfo) FromString(info string) error {
	return json.Unmarshal([]byte(info), sc)
}
