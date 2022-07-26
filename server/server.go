package server

import (
	"context"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/wsx864321/sweet_rpc/middleware"
)

type Server interface {
	Start()
	Close()
	RegisterService(serName string)
	RegisterMiddleware(middlewareFuncs ...middleware.MiddlewareFunc)
}

type Handle func(ctx context.Context) error

type server struct {
	opts        *Options
	middlewares []middleware.MiddlewareFunc
	serviceMap  map[string]*service
}

func NewServer(opts ...Option) *server {
	return &server{
		opts:        NewOptions(opts...),
		middlewares: make([]middleware.MiddlewareFunc, 0),
	}
}

// RegisterMiddleware 注册中间件
func (s *server) RegisterMiddleware(middlewareFuncs ...middleware.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewareFuncs...)
}

func (s *server) Start() {
	// 获取listener
	// logic

	// 注册服务
	for name, _ := range s.serviceMap {
		service := &discov.Service{
			Name: name,
			Endpoints: []*discov.Endpoint{
				{
					ServiceName: name,
					IP:          s.opts.IP,
					Port:        s.opts.Port,
					Enable:      true,
				},
			},
		}
		s.opts.Discovery.Register(context.TODO(), service)
	}
}
