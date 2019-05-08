package client

// Token 验证token并得到一个用户的永久唯一的UID
//       此UID应是客户端在各个服务器中的唯一标识
type Token interface {
	// GetSessionID 根据token返回一个临时的sessionid
	//              此sessionid应有有效期，代理服务器会将此sessionid返回客户端
	//       string 客户端发来的token
	GetSessionID(string) (uint32, error)
	// GetUID 根据客户端发来的sessionid得到用户的真实UID
	GetUID(uint32) (uint32, error)
}
