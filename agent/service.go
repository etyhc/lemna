package agent

import "lemna/utils"

var ServiceID uint32

func init() {
	ServiceID = utils.ID32Gen()
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

func (s *Service) Run() error {
	err := s.sp.Run()
	if err != nil {
		return err
	}
	return s.cp.Run()
}
