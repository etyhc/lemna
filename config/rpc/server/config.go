package server

import "encoding/json"

type Config struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Type int32  `json:"type"`
}

func (s *Config) String() string {
	ret, _ := json.Marshal(*s)
	return string(ret)
}

func (s *Config) Init(info string) error {
	return json.Unmarshal([]byte(info), s)
}
