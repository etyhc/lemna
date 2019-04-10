package impl

import fmt "fmt"

type SimpleToken struct {
	data map[string]int32
}

var tokenMap = map[string]int32{
	"token1": 1,
	"token2": 2}

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
