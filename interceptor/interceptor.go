package interceptor

import (
	"context"
)

type ServerInterceptor func(ctx context.Context, req interface{}, h Handler) (interface{}, error)

type Handler func(ctx context.Context, req interface{}) (interface{}, error)

type ClientInterceptor func(ctx context.Context, req, resp interface{}, h Invoker) error

type Invoker func(ctx context.Context, req, resp interface{}) error

// ServerIntercept server端拦截器
func ServerIntercept(ctx context.Context, req interface{}, ceps []ServerInterceptor, h Handler) (interface{}, error) {
	if len(ceps) == 0 {
		return h(ctx, req)
	}
	return ceps[0](ctx, req, getHandler(0, ceps, h))
}

func getHandler(cur int, ceps []ServerInterceptor, h Handler) Handler {
	if cur == len(ceps)-1 {
		return h
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return ceps[cur+1](ctx, req, getHandler(cur+1, ceps, h))
	}
}

// ClientIntercept client端拦截器
func ClientIntercept(ctx context.Context, req, resp interface{}, ceps []ClientInterceptor, i Invoker) error {
	if len(ceps) == 0 {
		return i(ctx, req, resp)
	}
	return ceps[0](ctx, req, resp, getInvoker(0, ceps, i))
}

func getInvoker(cur int, ceps []ClientInterceptor, i Invoker) Invoker {
	if cur == len(ceps)-1 {
		return i
	}
	return func(ctx context.Context, req, resp interface{}) error {
		return ceps[cur+1](ctx, req, resp, getInvoker(cur+1, ceps, i))
	}
}
