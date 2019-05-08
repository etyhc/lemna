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
	sid     uint32
	uid     uint32
	timeout int32
}

var tokenDB = map[string]*token{
	"token1": {1, 1, 0},
	"token2": {2, 2, 0}}

// GetUID  接口实现
func (st *SimpleToken) GetUID(sessionid uint32) (uint32, error) {
	for _, tk := range st.db {
		if tk.sid == sessionid {
			if atomic.LoadInt32(&tk.timeout) > 0 {
				return tk.uid, nil
			}
			return 0, fmt.Errorf("<sessionid=%d> timeout", sessionid)
		}
	}
	return 0, fmt.Errorf("no UID with<sessionid=%d>", sessionid)
}

// GetSessionID 接口实现
func (st *SimpleToken) GetSessionID(token string) (uint32, error) {
	tk, ok := st.db[token]
	if ok {
		atomic.StoreInt32(&tk.timeout, 5)
		return tk.sid, nil
	}
	return 0, fmt.Errorf("invaild SimpleToken %s", token)
}

// NewSimpleToken 新的Token服务
func NewSimpleToken() (st *SimpleToken) {
	st = &SimpleToken{}
	st.db = tokenDB
	go func() {
		tick := time.NewTicker(time.Duration(time.Second))
		defer tick.Stop()
		for {
			<-tick.C
			for _, tk := range st.db {
				if atomic.LoadInt32(&tk.timeout) > 0 {
				}
			}
		}
	}()
	return
}
