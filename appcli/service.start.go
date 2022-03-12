package appcli

import (
	"github.com/kardianos/service"
	"github.com/zhiyunliu/velocity/log"
)

//Start Start
func (p *ServiceApp) Start(s service.Service) (err error) {
	log.Infof("服务启动:%s", p.manager.Name())
	return p.run()
}
