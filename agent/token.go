package agent

// Token 接口，根据token返回一个用户的sessionid
//       此id应是客户端在各个服务器中的唯一标识
type Token interface {
	GetSessionID(string) (int32, error)
}
