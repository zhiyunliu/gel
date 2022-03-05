/*
 * @Author: lwnmengjing
 * @Date: 2021/6/7 5:39 下午
 * @Last Modified by: lwnmengjing
 * @Last Modified time: 2021/6/7 5:39 下午
 */

package server

import (
	"github.com/zhiyunliu/velocity/libs/types"
)

type Manager interface {
	Name() string
	Add(...Runnable)
	Start() error
}

type Runnable interface {
	Name() string
	// Start 启动
	Start() error

	// Stop 关闭
	Stop() error

	Status() string

	// Attempt 是否允许启动
	Attempt() bool

	Port() uint64

	Metadata() types.XMap

	//String 格式化
	String() string
}
