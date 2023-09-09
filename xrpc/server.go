package xrpc

import (
	"context"
	"fmt"

	"github.com/zhiyunliu/glue/config"
)

type Server interface {
	GetProto() string
	GetAddr() string
	Serve(ctx context.Context) (err error)
	Stop(ctx context.Context) error
}

// ServerResover 定义配置文件转换方法
type ServerResover interface {
	Name() string
	Resolve(name string, setting config.Config) (Server, error)
}

var serverResolvers = make(map[string]ServerResover)

// Register 注册配置文件适配器
func RegisterServer(resolver ServerResover) {
	proto := resolver.Name()
	if _, ok := serverResolvers[proto]; ok {
		panic(fmt.Errorf("xrpc: 不能重复注册:%s", proto))
	}
	serverResolvers[proto] = resolver
}

// Deregister 清理配置适配器
func DeregisterServer(name string) {
	delete(serverResolvers, name)
}

// NewServer 根据适配器名称及参数返回配置处理器
func NewServer(proto string, setting config.Config) (Server, error) {
	resolver, ok := serverResolvers[proto]
	if !ok {
		return nil, fmt.Errorf("xrpc: 未知的协议类型:%s", proto)
	}
	return resolver.Resolve(proto, setting)
}
