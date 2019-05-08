package utils

import (
	"fmt"
	"net"
)

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
