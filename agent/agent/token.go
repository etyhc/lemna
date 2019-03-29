package main

import fmt "fmt"

type token struct {
	data map[string]int32
}

var tokenMap = map[string]int32{
	"token1": 1,
	"token2": 2}

func (t *token) GetSessionID(token string) (id int32, err error) {
	var ok bool
	id, ok = t.data[token]
	if !ok {
		err = fmt.Errorf("invaild token %s", token)
	}
	return
}
