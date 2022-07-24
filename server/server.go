package server

import "github.com/wsx864321/sweet_rpc/middleware"

type Server struct {
	opts        *Options
	middlewares []middleware.MiddlewareFunc
}

func NewServer(opts ...Option) *Server {
	return &Server{
		opts:        NewOptions(opts...),
		middlewares: make([]middleware.MiddlewareFunc, 0),
	}
}

// RegisterMiddleware 注册中间件
func (s *Server) RegisterMiddleware(middlewareFuncs ...middleware.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewareFuncs...)
}
