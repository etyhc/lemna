package proto

//Builder protobuf消息构建器
type Builder struct {
}

//InvalidTarget 构建无效转发目标消息
func (b Builder) InvalidTarget() interface{} {
	return InvalidTargetMsg{}
}
