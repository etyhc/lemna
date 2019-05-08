package utils

import (
	"net"
	"time"
)

// ID32GenWithSalt 根据机器mac地址和salt生成32位ID
//                 当机器和salt不变,ID不变，保证一致性
func ID32GenWithSalt(salt string) uint32 {
	ifs, _ := net.Interfaces()
	hash := HashFnv1a(salt)
	for _, i := range ifs {
		hash = hash ^ HashFnv1a(i.HardwareAddr.String())
	}
	return hash
}

// ID32Gen 根据当前时间和机器mac地址生成32位ID，不保证一致性
func ID32Gen() uint32 {
	ifs, _ := net.Interfaces()
	var hash uint32
	for _, i := range ifs {
		hash = hash ^ HashFnv1a(i.HardwareAddr.String())
	}
	n := time.Now().UnixNano()
	return uint32(n>>32) ^ uint32(n) ^ hash

}
