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
	stp TargetPool
	ctp TargetPool
}

// NewService 新建代理服务
//         sp 服务器池
//         cp 客户端池
func NewService(stp, ctp TargetPool) *Service {
	return &Service{stp: stp, ctp: ctp}
}

// Run 非阻塞运行服务器池，阻塞运行客户端池
func (s *Service) Run() error {
	go func() {
		err := s.stp.Run(s.ctp)
		logger.Error(err)
	}()
	return s.ctp.Run(s.stp)
}
