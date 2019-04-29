package rpc

import (
	"lemna/content"
	"lemna/logger"
	"reflect"

	"github.com/gomodule/redigo/redis"
)

// Channel redis 频道
type Channel struct {
	Addr string //频道地址
}

// Publish content.Channel的redis实现
func (c *Channel) Publish(ctt content.Content) error {
	rc, err := redis.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	rc.Do("PUBLISH", ctt.Topic(), ctt.ToString())
	rc.Close()
	return nil
}

// Subscribe content.Channel的redis实现
//           redis无法订阅到订阅之前发布的消息
func (c *Channel) Subscribe(ctt content.Content) (<-chan content.Content, error) {
	rc, err := redis.Dial("tcp", c.Addr)
	if err != nil {
		return nil, err
	}
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.Subscribe(redis.Args{}.AddFlat(ctt.Topic())); err != nil {
		rc.Close()
		return nil, err
	}

	ret := make(chan content.Content)
	go func() {
		defer rc.Close()
		for {
			switch n := psc.Receive().(type) {
			case error:
				return
			case redis.Message:
				if n.Channel == ctt.Topic() {
					content := reflect.New(reflect.TypeOf(ctt).Elem()).Interface().(content.Content)
					err = content.FromString(string(n.Data))
					if err == nil {
						logger.Debug(content)
						ret <- content
					} else {
						logger.Error(err)
					}
				}
			default:
			}
		}
	}()
	return ret, nil
}
