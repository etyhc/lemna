package utils

import (
	"fmt"
	"net"
)

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
