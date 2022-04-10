package mqc

import (
	"bytes"
	"context"
	sctx "context"
	"encoding/json"
	"io"

	"github.com/zhiyunliu/gel/contrib/alloter"
	"github.com/zhiyunliu/gel/queue"
)

const XRemoteHeader = "x-remote-addr"

var _ alloter.IRequest = (*Request)(nil)

//Request 处理任务请求
type Request struct {
	ctx  sctx.Context
	task *Task
	queue.IMQCMessage
	method string
	params map[string]string
	header map[string]string
	body   cbody //map[string]string
}

//NewRequest 构建任务请求
func NewRequest(task *Task, m queue.IMQCMessage) (r *Request, err error) {

	r = &Request{

		IMQCMessage: m,
		task:        task,
		method:      "GET",
		params:      make(map[string]string),
	}

	//将消息原串转换为map
	message := m.GetMessage()

	r.header = message.Header()
	r.body = message.Body()
	r.ctx = context.Background()

	if r.header["Content-Type"] == "" {
		r.header["Content-Type"] = "application/json"
	}

	return r, nil
}

//GetName 获取任务名称
func (m *Request) GetName() string {
	return m.task.Queue
}

//GetService 服务名
func (m *Request) GetService() string {
	return m.task.GetService()
}

//GetMethod 方法名
func (m *Request) GetMethod() string {
	return m.method
}

func (m *Request) Params() map[string]string {
	return m.params
}

func (m *Request) GetHeader() map[string]string {
	return m.header
}

func (m *Request) Body() []byte {
	bytes, _ := json.Marshal(m.body)
	return bytes
}

func (m *Request) GetRemoteAddr() string {
	return m.header[XRemoteHeader]
}

func (m *Request) Context() sctx.Context {
	return m.ctx
}
func (m *Request) WithContext(ctx sctx.Context) alloter.IRequest {
	m.ctx = ctx
	return m
}

type Body interface {
	io.Reader
	Scan(obj interface{}) error
}

type cbody map[string]interface{}

func (b cbody) Read(p []byte) (n int, err error) {
	bodyBytes, err := json.Marshal(b)
	if err != nil {
		return 0, err
	}
	return bytes.NewReader(bodyBytes).Read(p)
}

func (b cbody) Scan(obj interface{}) error {
	bytes, err := json.Marshal(b)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, obj)
}
