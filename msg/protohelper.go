package msg

import (
	"lemna/utils"
	"reflect"

	proto "github.com/golang/protobuf/proto"
)

//ProtoInfo Protobuf基本信息
type ProtoInfo struct {
	id   uint32
	name string
	elem reflect.Type
}

//ID BaseInfo 的Protobuf实现
func (pi ProtoInfo) ID() uint32 {
	return pi.id
}

//Name BaseInfo 的Protobuf实现
func (pi ProtoInfo) Name() string {
	return pi.name
}

//ProtoHelper Protobuf消息辅助类
type ProtoHelper struct {
}

//Extract Helper的Protobuf实现
func (ph ProtoHelper) Extract(msg interface{}) Info {
	elem := reflect.TypeOf(msg).Elem()
	return ProtoInfo{id: utils.HashFnv1a(elem.Name()), name: elem.Name(), elem: elem}
}

//ToRaw Helper的Protobuf实现
func (ph ProtoHelper) ToRaw(msg interface{}) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

//FromRaw Helper的Protobuf实现
func (ph ProtoHelper) FromRaw(info Info, raw []byte) (interface{}, error) {
	msg := reflect.New(info.(ProtoInfo).elem).Interface()
	err := proto.Unmarshal(raw, msg.(proto.Message))
	return msg, err
}
