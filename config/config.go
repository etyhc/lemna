package config

// Stringer 配置的序列化反序列化接口
type Stringer interface {
	String() string
	FromString(string) error
}

// Channel 配置频道，可以在此发布和订阅配置
//         目前实现了一个基于grpc的频道服务器
type Channel interface {
	Publish(string, Stringer) error
	Subscribe(string, Stringer) (<-chan Stringer, error)
}
