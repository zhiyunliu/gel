package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kardianos/service"
	"github.com/urfave/cli"
	"github.com/zhiyunliu/velocity/cli/server"
	"github.com/zhiyunliu/velocity/configs"
	"github.com/zhiyunliu/velocity/libs"
)

type AppService struct {
	service.Service
	ServiceName string
	DisplayName string
	Description string
	Arguments   []string
}

//GetService GetService
func GetService(c *cli.Context, args ...string) (velocitySrv *AppService, err error) {
	srvApp := GetSrvApp(c)
	//1. 构建服务配置
	cfg := GetSrvConfig(srvApp.config, args...)
	//2.创建本地服务
	appSrv, err := service.New(srvApp, cfg)
	if err != nil {
		return nil, err
	}
	bytes, _ := json.Marshal(cfg)

	fmt.Println("sss:", string(bytes))
	return &AppService{
		Service:     appSrv,
		ServiceName: cfg.Name,
		DisplayName: cfg.DisplayName,
		Description: cfg.Description,
		Arguments:   cfg.Arguments,
	}, err
}

//GetSrvConfig SrvCfg
func GetSrvConfig(appCfg *configs.AppSetting, args ...string) *service.Config {
	path, _ := filepath.Abs(os.Args[0])

	svcName := fmt.Sprintf("%s_%s", appCfg.SysName, libs.Md5(path)[:8])

	bytes, _ := json.Marshal(appCfg)
	fmt.Println("appCfg:", path, string(bytes))

	cfg := &service.Config{
		Name:         svcName,
		DisplayName:  svcName,
		Description:  appCfg.Usage,
		Arguments:    args,
		Dependencies: []string{"After=network.target syslog.target"},
	}
	cfg.WorkingDirectory = filepath.Dir(path)
	// cfg.Option = make(map[string]interface{})
	// cfg.Option["LimitNOFILE"] = 10240
	return cfg
}

//GetSrvApp SrvCfg
func GetSrvApp(c *cli.Context) *ServiceApp {
	server := c.App.Metadata["server"].(server.Server)
	appCfg := c.App.Metadata["config"].(*configs.AppSetting)
	initAppConfig(appCfg)
	return &ServiceApp{
		c:      c,
		server: server,
		config: appCfg,
	}
}

//ServiceApp ServiceApp
type ServiceApp struct {
	c      *cli.Context
	server server.Server
	config *configs.AppSetting
	trace  itrace
}

func initAppConfig(config *configs.AppSetting) {
	if config.Addr == "" {
		config.Addr = ":8081"
	}

	if config.PlatName == "" {
		config.PlatName = "default"
	}
	if config.SysName == "" {
		config.SysName = filepath.Base(os.Args[0])
	}
	if config.ClusterName == "" {
		config.ClusterName = "prod"
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}
}
