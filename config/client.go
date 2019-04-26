package config

import "encoding/json"

// ClientSession 客户端session
//               客户端验证成功后，发布自己的状态
type ClientSession struct {
	SessionID int32  `json:"id"` //客户端在代理服务器的唯一ID
	AgentID   int32  `json:"aid"`
	UID       string `json:"uid"`
}

func (sc *ClientSession) String() string {
	ret, _ := json.Marshal(*sc)
	return string(ret)
}

// FromString 从字串中初始化服务器配置,使用json编码
func (sc *ClientSession) FromString(info string) error {
	return json.Unmarshal([]byte(info), sc)
}
