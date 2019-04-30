package redis

import (
	"lemna/content"
	"lemna/logger"
	"reflect"

	"github.com/gomodule/redigo/redis"
)

const REDISADDR = ":6379"

// Channel redis 频道
type Channel struct {
	Addr string //频道地址
}

// Publish content.Channel的redis实现
func (c *Channel) Publish(ctt content.Content) error {
	cts, err := content.ToJSON(ctt)
	if err != nil {
		return nil
	}
	rc, err := redis.Dial("tcp", c.Addr)
	if err != nil {
		return err
	}
	_, err = rc.Do("PUBLISH", ctt.Topic(), cts)
	//logger.Debugf("pub:<%s>%s rep:%d", ctt.Topic(), string(cts), reply)
	rc.Close()
	return err
}

// Subscribe content.Channel的redis实现
//           redis无法订阅到订阅之前发布的消息
func (c *Channel) Subscribe(ctt content.Content) (<-chan content.Content, error) {
	rc, err := redis.Dial("tcp", c.Addr)
	if err != nil {
		return nil, err
	}
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.Subscribe(ctt.Topic()); err != nil {
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
				//logger.Debug("sub:", n, ctt)
				if n.Channel == ctt.Topic() {
					c := reflect.New(reflect.TypeOf(ctt).Elem()).Interface().(content.Content)
					c, err := content.FromJSON(c, n.Data)
					if err == nil {
						//logger.Debug(c)
						ret <- c
					} else {
						logger.Error(err)
					}
				}
			default:
				//logger.Debug("def:", n)
			}
		}
	}()
	return ret, nil
}
