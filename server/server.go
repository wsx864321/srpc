package server

import (
	"context"
	"fmt"
	"github.com/wsx864321/sweet_rpc/codec"
	"github.com/wsx864321/sweet_rpc/transport"
	"net"
	"time"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/wsx864321/sweet_rpc/interceptor"
)

type Server interface {
	Start()
	Close()
	RegisterService(serName string)
	RegisterMiddleware(middlewareFuncs ...interceptor.MiddlewareFunc)
}

type Handle func(ctx context.Context) error

type server struct {
	opts         *Options
	codec        codec.Codec
	ctx          context.Context    // Each service is managed in one context
	cancel       context.CancelFunc // controller of context
	middlewares  []interceptor.MiddlewareFunc
	serviceMap   map[string]*service
	transportMgr *transport.ServerTransportMgr
}

// NewServer 生成一个server
func NewServer(opts ...Option) *server {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		opts:         NewOptions(opts...),
		ctx:          ctx,
		cancel:       cancel,
		middlewares:  make([]interceptor.MiddlewareFunc, 0),
		transportMgr: transport.NewServerTransportMgr(),
	}
}

// RegisterMiddleware 注册中间件
func (s *server) RegisterMiddleware(middlewareFuncs ...interceptor.MiddlewareFunc) {
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
			// 开启一个协程进行rpc
			go s.process(accept)
		}

	}

}

// process logic
func (s *server) process(conn net.Conn) {
	// 1.提取数据
	//msg, err := s.extractMessage(conn)
	//if err != nil {
	//	s.opts.Logger.Errorf(context.TODO(), "extractMessage error:%v", err.Error())
	//	return
	//}

	// 2.执行中间件

}

// extractMessage 提取message内容
func (s *server) extractMessage(conn net.Conn) (*codec.Message, error) {
	// 1.设置读取超时时间
	if err := conn.SetReadDeadline(time.Now().Add(s.opts.ReadTimeout)); err != nil {
		return nil, err
	}

	// 2.读取头部内容
	headerData := make([]byte, s.codec.GetHeaderLength())
	if err := s.readFixedData(conn, headerData); err != nil {
		return nil, err
	}

	header, err := s.codec.DecodeHeader(headerData)
	if err != nil {
		return nil, err
	}

	// 3.读取message内容
	body := make([]byte, s.codec.GetBodyLength(header))
	if err = s.readFixedData(conn, body); err != nil {
		return nil, err
	}

	return s.codec.DecodeBody(header, body)
}

// readFixedData 读取固定长度内容
func (s *server) readFixedData(conn net.Conn, buf []byte) error {
	var (
		pos       = 0
		totalSize = len(buf)
	)
	for {
		c, err := conn.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + c
		if pos == totalSize {
			break
		}
	}
	return nil
}
