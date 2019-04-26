package impl

import fmt "fmt"

type SimpleToken struct {
	data map[string]int32
}

var tokenMap = map[string]int32{
	"token1": 1,
	"token2": 2}

func (t *SimpleToken) GetUID(sessionid int32) (int32, error) {
	for _, v := range t.data {
		if v == sessionid {
			return sessionid, nil
		}
	}
	return 0, fmt.Errorf("no UID with<sessionid=%d>", sessionid)
}

func (t *SimpleToken) GetSessionID(token string) (int32, error) {
	id, ok := t.data[token]
	if ok {
		return id, nil
	}
	return id, fmt.Errorf("invaild SimpleToken %s", token)
}

func NewSimpleToken() (t *SimpleToken) {
	t = &SimpleToken{}
	t.data = tokenMap
	return
}
