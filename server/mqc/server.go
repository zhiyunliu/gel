package mqc

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/zhiyunliu/glue/config"
	_ "github.com/zhiyunliu/glue/contrib/xmqc/alloter"
	"github.com/zhiyunliu/glue/engine"
	"github.com/zhiyunliu/glue/global"
	"github.com/zhiyunliu/glue/log"
	"github.com/zhiyunliu/glue/middleware"
	"github.com/zhiyunliu/glue/server"
	"github.com/zhiyunliu/glue/transport"
	"github.com/zhiyunliu/glue/xmqc"
)

type Server struct {
	name    string
	server  xmqc.Server
	ctx     context.Context
	opts    *options
	started bool
}

var _ transport.Server = (*Server)(nil)

// New 实例化
func New(name string, opts ...Option) *Server {
	s := &Server{
		name: name,
		opts: setDefaultOption(),
	}
	s.Options(opts...)

	return s
}

// Options 设置参数
func (e *Server) Options(opts ...Option) {
	for _, o := range opts {
		o(e.opts)
	}
}

func (e *Server) Name() string {
	if e.name == "" {
		e.name = e.Type()
	}
	return e.name
}

func (e *Server) Type() string {
	return Type
}

// ServiceName 服务名称
func (s *Server) ServiceName() string {
	return s.opts.serviceName
}

func (e *Server) Endpoint() *url.URL {
	return transport.NewEndpoint(e.Type(), fmt.Sprintf("%s:%d", global.LocalIp, 1987))
}

func (e *Server) Config(cfg config.Config) {
	if cfg == nil {
		return
	}
	e.Options(WithConfig(cfg))
	cfg.Get(e.serverPath()).Scan(e.opts.srvCfg)
}

// Start 开始
func (e *Server) Start(ctx context.Context) (err error) {
	if e.opts.srvCfg.Config.Status == server.StatusStop {
		return nil
	}

	e.ctx = transport.WithServerContext(ctx, e)

	for _, m := range e.opts.srvCfg.Middlewares {
		e.opts.router.Use(middleware.Resolve(&m))
	}

	e.server, err = xmqc.NewServer(e.opts.srvCfg.Config.Proto,
		e.opts.router,
		e.opts.config.Get(e.serverPath()),
		engine.WithConfig(e.opts.config),
		engine.WithLogOptions(e.opts.logOpts),
		engine.WithSrvType(e.Type()),
		engine.WithSrvName(e.Name()),
		engine.WithErrorEncoder(e.opts.encErr),
		engine.WithRequestDecoder(e.opts.decReq),
		engine.WithResponseEncoder(e.opts.encResp),
	)

	if err != nil {
		return
	}

	errChan := make(chan error, 1)
	log.Infof("MQC Server [%s] listening on %s", e.name, e.opts.srvCfg.Config.String())
	done := make(chan struct{})

	go func() {
		e.started = true
		errChan <- e.server.Serve(e.ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		errChan <- nil
	}
	err = <-errChan
	if err != nil {
		log.Errorf("MQC Server [%s] start error: %s", e.name, err.Error())
		return err
	}

	if len(e.opts.startedHooks) > 0 {
		for _, fn := range e.opts.startedHooks {
			err := fn(ctx)
			if err != nil {
				log.Errorf("MQC Server [%s] StartedHooks:%+v", e.name, err)
				return err
			}
		}
	}
	log.Infof("MQC Server [%s] start completed", e.name)
	return nil
}

// Attempt 判断是否可以启动
func (e *Server) Attempt() bool {
	return !e.started
}

// Shutdown 停止
func (e *Server) Stop(ctx context.Context) error {
	err := e.server.Stop(ctx)
	if err != nil {
		log.Errorf("MQC Server [%s] stop error: %s", e.name, err.Error())
		return err
	}

	if len(e.opts.endHooks) > 0 {
		for _, fn := range e.opts.endHooks {
			err := fn(ctx)
			if err != nil {
				log.Errorf("MQC Server [%s] EndHook:", e.name, err)
				return err
			}
		}
	}
	log.Infof("MQC Server [%s] stop completed", e.name)
	return nil

}

func (e *Server) Use(middlewares ...middleware.Middleware) {
	e.opts.router.Use(middlewares...)
}

func (e *Server) Group(group string, middlewares ...middleware.Middleware) *engine.RouterGroup {
	return e.opts.router.Group(group, middlewares...)
}

func (e *Server) Handle(queue string, obj interface{}) {
	e.opts.router.Handle(getService(queue), obj, engine.MethodGet)
}

func (e *Server) serverPath() string {
	return fmt.Sprintf("servers.%s", e.Name())
}
