package msg

import (
	"lemna/arpc"
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

//ID Info的Protobuf实现
func (pi ProtoInfo) ID() uint32 {
	return pi.id
}

//Name Info的Protobuf实现
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
func (ph ProtoHelper) ToRaw(msg interface{}) (uint32, []byte, error) {
	id := ph.Extract(msg).ID()
	raw, err := proto.Marshal(msg.(proto.Message))
	return id, raw, err
}

//FromRaw Helper的Protobuf实现
func (ph ProtoHelper) FromRaw(info Info, raw []byte) (interface{}, error) {
	msg := reflect.New(info.(ProtoInfo).elem).Interface()
	err := proto.Unmarshal(raw, msg.(proto.Message))
	return msg, err
}

//Wrap Helper的Protobuf实现
func (ph ProtoHelper) Wrap(target uint32, msg interface{}) (ForwardMsg, error) {
	id, rm, err := ph.ToRaw(msg)
	if err != nil {
		return nil, err
	}
	return &arpc.ForwardMsg{Target: target, Mid: id, Raw: rm}, nil
}

//WrapMM Helper的Protobuf实现
func (ph ProtoHelper) WrapMM(targets []uint32, msg interface{}) (MulticastMsg, error) {
	id, rm, err := ph.ToRaw(msg)
	if err != nil {
		return nil, err
	}
	return &arpc.MulticastMsg{Targets: targets, Mid: id, Raw: rm}, nil
}
