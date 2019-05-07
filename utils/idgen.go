package utils

import (
	"net"
	"time"
)

func ID32GenWithSalt(salt string) uint32 {
	ifs, _ := net.Interfaces()
	hash := HashFnv1a(salt)
	for _, i := range ifs {
		hash = hash ^ HashFnv1a(i.HardwareAddr.String())
	}
	return hash
}

func ID32Gen() uint32 {
	ifs, _ := net.Interfaces()
	var hash uint32
	for _, i := range ifs {
		hash = hash ^ HashFnv1a(i.HardwareAddr.String())
	}
	n := time.Now().UnixNano()
	return uint32(n>>32) ^ uint32(n) ^ hash

}
