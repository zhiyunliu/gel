package grpc

import (
	"fmt"
	"net/http"

	"context"

	"github.com/zhiyunliu/glue/contrib/alloter"
	"github.com/zhiyunliu/glue/contrib/xrpc/grpc/grpcproto"
	"github.com/zhiyunliu/golibs/bytesconv"
)

// processor cron管理程序，用于管理多个任务的执行，暂停，恢复，动态添加，移除
type processor struct {
	engine *alloter.Engine
}

// NewProcessor 创建processor
func newProcessor(engine *alloter.Engine) (p *processor, err error) {
	p = &processor{}
	p.engine = engine
	return p, nil
}

func (s *processor) Process(ctx context.Context, request *grpcproto.Request) (response *grpcproto.Response, err error) {
	response = &grpcproto.Response{}
	//转换输入参数

	req, err := newServerRequest(ctx, request)
	if err != nil {
		response.Status = int32(http.StatusNotAcceptable)
		response.Result = bytesconv.StringToBytes(fmt.Sprintf("输入参数有误:%v", err))
		return response, nil
	}

	resp := newServerResponse()

	//发起本地处理
	err = s.engine.HandleRequest(req, resp)
	if err != nil {
		response.Status = int32(http.StatusInternalServerError)
		response.Result = bytesconv.StringToBytes(fmt.Sprintf("处理请求有误%s", err.Error()))
		return response, nil
	}

	//处理响应内容
	response.Status = int32(resp.Status())
	response.Result = resp.buffer.Bytes()
	response.Headers = resp.Header()
	return response, nil

}
