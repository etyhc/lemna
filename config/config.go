//Package config 配置订阅服务
package config

// Stringer 配置的序列化反序列化接口
type Stringer interface {
	// String 将配置序列化为字串
	String() string
	// FromString 从字符串初始化配置,失败返回error
	FromString(string) error
}

// Channel 配置频道，可以在此发布和订阅配置
type Channel interface {
	//  Publish 发布配置
	//   string 配置名字
	// Stringer 配置接口
	Publish(string, Stringer) error
	//     Subscribe 订阅配置
	//        string 订阅配置名字
	//      Stringer 配置接口
	// chan Stringer 返回一个配置chan，在这里读取订阅的配置
	Subscribe(string, Stringer) (<-chan Stringer, error)
}
