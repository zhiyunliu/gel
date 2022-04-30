package circuitbreaker

import (
	"github.com/zhiyunliu/gel/context"

	"github.com/go-kratos/aegis/circuitbreaker"
	"github.com/go-kratos/aegis/circuitbreaker/sre"
	"github.com/zhiyunliu/gel/errors"
	"github.com/zhiyunliu/gel/golibs/group"
	"github.com/zhiyunliu/gel/middleware"
	"github.com/zhiyunliu/gel/transport"
)

// ErrNotAllowed is request failed due to circuit breaker triggered.
var ErrNotAllowed = errors.New(503, "request failed due to circuit breaker triggered")

// Option is circuit breaker option.
type Option func(*options)

// WithGroup with circuit breaker group.
// NOTE: implements generics circuitbreaker.CircuitBreaker
func WithGroup(g *group.Group) Option {
	return func(o *options) {
		o.group = g
	}
}

type options struct {
	group *group.Group
}

// Client circuitbreaker middleware will return errBreakerTriggered when the circuit
// breaker is triggered and the request is rejected directly.
func Client(opts ...Option) middleware.Middleware {
	opt := &options{
		group: group.NewGroup(func() interface{} {
			return sre.NewBreaker()
		}),
	}
	for _, o := range opts {
		o(opt)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context) (reply interface{}) {
			orgCtx := ctx.Context()
			info, _ := transport.FromClientContext(orgCtx)
			breaker := opt.group.Get(info.Operation()).(circuitbreaker.CircuitBreaker)
			if err := breaker.Allow(); err != nil {
				// rejected
				// NOTE: when client reject requets locally,
				// continue add counter let the drop ratio higher.
				breaker.MarkFailed()
				return ErrNotAllowed
			}
			// allowed
			reply = handler(ctx)
			var err error
			if rerr, ok := reply.(error); ok {
				err = rerr
			}

			if err != nil && errors.IsInternalServer(err) {
				breaker.MarkFailed()
			} else {
				breaker.MarkSuccess()
			}
			return reply
		}
	}
}