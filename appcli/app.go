package appcli

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"
	"github.com/zhiyunliu/velocity/appcli/keys"
	"github.com/zhiyunliu/velocity/global"
)

//App  cli app
type App struct {
	cliApp  *cli.App
	options *Options
}

//Start 启动应用程序
func (a *App) Start() error {
	a.cliApp.Commands = GetCmds(a.options)
	return a.cliApp.Run(os.Args)
}

//New 创建app
func New(opts ...Option) *App {

	app := &App{options: &Options{}}
	for _, opt := range opts {
		opt(app.options)
	}

	app.cliApp = cli.NewApp()
	app.cliApp.Name = filepath.Base(os.Args[0])
	app.cliApp.Version = global.Version
	app.cliApp.Usage = global.Usage
	cli.HelpFlag = cli.BoolFlag{
		Name:  "help,h",
		Usage: "查看帮助信息",
	}
	cli.VersionFlag = cli.BoolFlag{
		Name:  "version,v",
		Usage: "查看版本信息",
	}
	app.App.Metadata = map[string]interface{}{}
	app.App.Metadata[keys.ManagerKey] = manager
	app.App.Metadata[keys.OptionsKey] = app.options

	return app
}
