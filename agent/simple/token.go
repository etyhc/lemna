package simple

import fmt "fmt"

type SimpleToken struct {
	data map[string]int32
}

var tokenMap = map[string]int32{
	"token1": 1,
	"token2": 2}

func (t *SimpleToken) GetSessionID(token string) (id int32, err error) {
	var ok bool
	id, ok = t.data[token]
	if !ok {
		err = fmt.Errorf("invaild SimpleToken %s", token)
	}
	return
}

func NewSimpleToken() (t *SimpleToken) {
	t = &SimpleToken{}
	t.data = tokenMap
	return
}
