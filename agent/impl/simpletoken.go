package impl

import (
	fmt "fmt"
	"time"
)

type SimpleToken struct {
	db map[string]token
}

type token struct {
	sid     int32
	timeout int32
}

var tokenDB = map[string]token{
	"token1": {1, 0},
	"token2": {2, 0}}

func (st *SimpleToken) GetUID(sessionid int32) (int32, error) {
	for _, tk := range st.db {
		if tk.sid == sessionid && tk.timeout > 0 {
			return tk.sid, nil
		}
	}
	return 0, fmt.Errorf("no UID with<sessionid=%d>", sessionid)
}

func (st *SimpleToken) GetSessionID(token string) (int32, error) {
	tk, ok := st.db[token]
	if ok {
		tk.timeout = 5
		return tk.sid, nil
	}
	return 0, fmt.Errorf("invaild SimpleToken %s", token)
}

func NewSimpleToken() (st *SimpleToken) {
	st = &SimpleToken{}
	st.db = tokenDB
	go func() {
		tick := time.NewTicker(time.Duration(time.Second))
		for {
			<-tick.C
			for _, tk := range st.db {
				if tk.timeout > 0 {
					tk.timeout--
				}
			}
		}
	}()
	return
}
