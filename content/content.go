//Package content 内容订阅服务，内容可以是任何东西
package content

import (
	"encoding/json"
	"reflect"
)

// Content 内容序列化反序列化接口
type Content interface {
	// Topic 内容主题
	Topic() string
}

// Channel 内容频道，可以在此发布和订阅内容
type Channel interface {
	// Publish 发布内容
	// Content 要发布的内容
	Publish(Content) error
	//    Subscribe 订阅内容
	//      Content 订阅的内容接口
	// chan Content 返回一个内容chan，在这里读取订阅的内容
	Subscribe(Content) (<-chan Content, error)
}

func Topic(v interface{}) string {
	return reflect.TypeOf(v).Elem().Name()
}

func ToJSON(ctt Content) ([]byte, error) {
	ret, err := json.Marshal(ctt)
	return ret, err
}

func FromJSON(ctt Content, jsonstr []byte) (Content, error) {
	err := json.Unmarshal(jsonstr, ctt)
	if err != nil {
		return nil, err
	}
	return ctt, nil
}
