package client

import (
	fmt "fmt"
	"sync/atomic"
	"time"
)

// SimpleToken 一个Token示例
type SimpleToken struct {
	db map[string]*token
}

type token struct {
	session string
	uid     uint32
	timeout int32
}

var tokenDB = map[string]*token{
	"token1": {"session1", 1, 0},
	"token2": {"session1", 2, 0}}

// GetUID  接口实现
func (st *SimpleToken) GetUID(session string) (uint32, error) {
	for _, tk := range st.db {
		if tk.session == session {
			if atomic.LoadInt32(&tk.timeout) > 0 {
				return tk.uid, nil
			}
			return 0, fmt.Errorf("<session=%s> timeout", session)
		}
	}
	return 0, fmt.Errorf("no UID with<sessionid=%s>", session)
}

// GetSession 接口实现
func (st *SimpleToken) GetSession(token string) (string, error) {
	tk, ok := st.db[token]
	if ok {
		atomic.StoreInt32(&tk.timeout, 5)
		return tk.session, nil
	}
	return "", fmt.Errorf("invaild SimpleToken %s", token)
}

// NewSimpleToken 新的Token服务
func NewSimpleToken() (st *SimpleToken) {
	st = &SimpleToken{}
	st.db = tokenDB
	go func() {
		tick := time.NewTicker(time.Second)
		defer tick.Stop()
		for {
			<-tick.C
			for _, tk := range st.db {
				if tk.timeout > 0 {
					atomic.AddInt32(&tk.timeout, -1)
				}
			}
		}
	}()
	return
}
