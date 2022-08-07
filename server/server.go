package server

import (
	"context"
	"fmt"
	"github.com/wsx864321/sweet_rpc/transport"
	"net"

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
	opts         *Options
	ctx          context.Context    // Each service is managed in one context
	cancel       context.CancelFunc // controller of context
	middlewares  []middleware.MiddlewareFunc
	serviceMap   map[string]*service
	transportMgr *transport.ServerTransportMgr
}

func NewServer(opts ...Option) *server {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		opts:         NewOptions(opts...),
		ctx:          ctx,
		cancel:       cancel,
		middlewares:  make([]middleware.MiddlewareFunc, 0),
		transportMgr: transport.NewServerTransportMgr(),
	}
}

// RegisterMiddleware 注册中间件
func (s *server) RegisterMiddleware(middlewareFuncs ...middleware.MiddlewareFunc) {
	s.middlewares = append(s.middlewares, middlewareFuncs...)
}

func (s *server) Start() {
	// 获取listener
	serverTransport := s.transportMgr.Gen(s.opts.Protocol)
	ln, err := serverTransport.Listen(fmt.Sprintf("%v:%v", s.opts.IP, s.opts.Port))
	if err != nil {
		panic(err)
	}

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

	go s.run(ln)
}

func (s *server) run(ln net.Listener) {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			accept, err := ln.Accept()
			if err != nil {
				s.opts.Logger.Errorf(s.ctx, "err:%v", err)
				continue
			}

			go s.process(accept)
		}

	}

}

func (s *server) process(accept net.Conn) {

}
