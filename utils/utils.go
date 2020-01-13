// Package utils 提供一些无法归类的功能.
package utils

import (
	"fmt"
	"hash/fnv"
	"net"
	"reflect"
	"time"
)

//IsNil 判断接口内容是否是空指针
func IsNil(i interface{}) bool {
	return (i == nil) || (reflect.ValueOf(i).Kind() == reflect.Ptr && reflect.ValueOf(i).IsNil())
}

// HashFnv1a 简化fvn1a算法调用
func HashFnv1a(str string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(str))
	return h.Sum32()
}

// PublishTCPAddr 过滤发布地址
//                发布地址无效返回""
//                发布地址指定，返回指定地址
//                发布地址未指定，找到第一个非loopback、非链路本地、单播地址返回
func PublishTCPAddr(addr string) string {
	tcpaddr, _ := net.ResolveTCPAddr("tcp", addr)
	if tcpaddr == nil {
		return ""
	}
	if tcpaddr.IP != nil && !tcpaddr.IP.IsUnspecified() {
		return tcpaddr.String()
	}

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		ip := addr.(*net.IPNet).IP
		if !ip.IsMulticast() &&
			!ip.IsLoopback() &&
			!ip.IsLinkLocalUnicast() {
			return fmt.Sprintf("%s:%d", ip.String(), tcpaddr.Port)
		}
	}
	return ""
}

// ID32GenWithSalt 根据机器mac地址和salt生成32位ID
//                 当机器和salt不变,ID不变，保证一致性
func ID32GenWithSalt(salt string) uint32 {
	return HashFnv1a(salt) ^ _machash
}

var _machash uint32

func init() {
	ifs, _ := net.Interfaces()
	for _, i := range ifs {
		_machash ^= HashFnv1a(i.HardwareAddr.String())
	}
}

// ID32Gen 根据当前时间和机器mac地址生成32位ID，不保证一致性
func ID32Gen() uint32 {
	n := time.Now().UnixNano()
	return uint32(n>>32) ^ uint32(n) ^ _machash

}
