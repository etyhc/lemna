package agent

import (
	"fmt"
	"reflect"

	"github.com/protobuf/proto"
)

type MsgHandler func(int32, interface{}, interface{})

type msgInfo struct {
	mType    reflect.Type
	mHandler MsgHandler
}

var msgMap = make(map[int32]msgInfo)

func MsgReg(id int32, msg interface{}, handler MsgHandler) {
	msgMap[id] = msgInfo{reflect.TypeOf(msg.(proto.Message)), handler}
}

func RawMsgEncode(t int32, msg proto.Message) (*RawMsg, error) {
	buf, err := proto.Marshal(msg)
	return &RawMsg{Type: t, Raw: buf}, err
}

func RawMsgHandle(id int32, t int32, raw []byte, stream interface{}) error {
	if info, ok := msgMap[t]; ok {
		msg := reflect.New(info.mType.Elem()).Interface()
		err := proto.Unmarshal(raw, msg.(proto.Message))
		if err == nil {
			info.mHandler(id, msg, stream)
		}
		return err
	}
	return fmt.Errorf("invalid message type %d", t)
}
