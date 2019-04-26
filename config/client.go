package config

import "encoding/json"

// ClientSession 客户端session
//               客户端验证成功后，发布自己的状态
type ClientSession struct {
	SessionID int32  `json:"sid"` //客户端在代理服务器上的唯一ID
	AgentID   int32  `json:"aid"` //客户端所在代理服务器唯一ID
	UID       string `json:"uid"` //客户端唯一ID
}

func (sc *ClientSession) String() string {
	ret, _ := json.Marshal(*sc)
	return string(ret)
}

// FromString 从字串中初始化服务器配置,使用json编码
func (sc *ClientSession) FromString(info string) error {
	return json.Unmarshal([]byte(info), sc)
}
