package agent

import (
	"lemna/utils"
	"net"
	"time"
)

var ServiceID uint32

func init() {
	ServiceID = genID()
}

type Service struct {
	sp TargetPool
	cp TargetPool
}

func NewService(sp, cp TargetPool) *Service {
	s := &Service{sp: sp, cp: cp}
	s.sp.Bind(s.cp)
	s.cp.Bind(s.sp)
	return s
}

func genID() uint32 {
	addrs, _ := net.InterfaceAddrs()
	var hash uint32
	for _, addr := range addrs {
		hash = hash ^ utils.HashFnv1a(addr.String())
	}
	return uint32(time.Now().Nanosecond()) ^ hash
}

func (s *Service) Run() error {
	err := s.sp.Run()
	if err != nil {
		return err
	}
	return s.cp.Run()
}
