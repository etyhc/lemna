package agent

import (
	"lemna/logger"
	"lemna/utils"
)

// ServiceID 代理服务ID，每次启动都会不一样,用于代理识别
var ServiceID uint32

func init() {
	ServiceID = utils.ID32Gen()
}

// Service 代理服务
type Service struct {
	sp TargetPool
	cp TargetPool
}

// NewService 新建代理服务
//         sp 服务器池
//         cp 客户端池
func NewService(sp, cp TargetPool) *Service {
	s := &Service{sp: sp, cp: cp}
	s.sp.Bind(s.cp)
	s.cp.Bind(s.sp)
	return s
}

// Run 非阻塞运行服务器池，阻塞运行客户端池
func (s *Service) Run() error {
	go func() {
		err := s.sp.Run()
		logger.Error(err)
	}()
	return s.cp.Run()
}
