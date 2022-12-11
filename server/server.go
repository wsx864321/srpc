package server

import (
	"context"
	"fmt"
	"github.com/wsx864321/sweet_rpc/codec"
	"github.com/wsx864321/sweet_rpc/codec/serialize"
	"github.com/wsx864321/sweet_rpc/transport"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/wsx864321/sweet_rpc/interceptor"
)

type handler func(ctx context.Context, body []byte) (interface{}, error)

type Server interface {
	Start()
	Stop()
	RegisterService(serName string, service interface{})
	RegisterMiddleware(interceptors ...interceptor.ServerInterceptor)
}

type server struct {
	opts         *Options
	codec        codec.Codec
	ctx          context.Context    // Each service is managed in one context
	cancel       context.CancelFunc // controller of context
	interceptors []interceptor.ServerInterceptor
	serviceMap   map[string]*service
	transportMgr *transport.ServerTransportMgr
}

type service struct {
	name    string
	methods map[string]handler
}

// NewServer 生成一个server
func NewServer(opts ...Option) *server {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		opts:         NewOptions(opts...),
		ctx:          ctx,
		cancel:       cancel,
		interceptors: make([]interceptor.ServerInterceptor, 0),
		transportMgr: transport.NewServerTransportMgr(),
	}
}

// RegisterMiddleware 注册中间件
func (s *server) RegisterMiddleware(interceptors ...interceptor.ServerInterceptor) {
	s.interceptors = append(s.interceptors, interceptors...)
}

// RegisterService 注册服务
func (s *server) RegisterService(serName string, srv interface{}) {
	svrType := reflect.TypeOf(srv)
	svrValue := reflect.ValueOf(srv)

	methods := make(map[string]handler)
	for i := 0; i < svrType.NumMethod(); i++ {
		method := svrType.Method(i)
		if err := s.checkMethod(method.Type); err != nil {
			panic(err)
		}

		methodHandler := func(ctx context.Context, body []byte) (interface{}, error) {
			req := reflect.New(method.Type.In(2).Elem()).Interface()
			if err := serialize.GetSerialize(s.opts.Serialize).Unmarshal(body, req); err != nil {
				return nil, err
			}

			h := func(ctx context.Context, req interface{}) (interface{}, error) {
				resp := svrType.Method(i).Func.Call([]reflect.Value{svrValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})

				if resp[0].IsValid() && resp[1].IsValid() {
					return nil, nil
				}

				if resp[0].IsValid() && !resp[1].IsValid() {
					return nil, resp[1].Interface().(error)
				}

				if !resp[0].IsValid() && resp[1].IsValid() {
					return resp[0].Interface(), nil
				}

				return resp[0].Interface(), resp[1].Interface().(error)
			}

			return interceptor.ServerIntercept(context.TODO(), req, s.interceptors, h)
		}
		methods[method.Name] = methodHandler
	}

	s.serviceMap[serName] = &service{
		serName,
		methods,
	}

	// 服务注册
	s.opts.Discovery.Register(context.TODO(), &discov.Service{
		Name: serName,
		Endpoints: []*discov.Endpoint{
			{
				ServiceName: serName,
				IP:          s.opts.IP,
				Port:        s.opts.Port,
				Enable:      true,
			},
		},
	})
}

func (s *server) checkMethod(methodType reflect.Type) error {
	if methodType.NumIn() != 3 {
		return fmt.Errorf("method %s invalid, the number of params != 2", methodType.Name())
	}

	if methodType.NumOut() != 2 {
		return fmt.Errorf("method %s invalid, the number of params != 2", methodType.Name())
	}

	var ctx *context.Context
	if !methodType.In(1).Implements(reflect.TypeOf(ctx).Elem()) {
		return fmt.Errorf("method %s invalid, first param is not context", methodType.Name())
	}

	if methodType.In(2).Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, second param is not ptr", methodType.Name())
	}

	if methodType.Out(0).Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, first reply type is not a pointer", methodType.Name())
	}

	var err *error
	if methodType.Out(1).Implements(reflect.TypeOf(err).Elem()) {
		return fmt.Errorf("method %s invalid, second reply is not error", methodType.Name())
	}

	return nil
}

// Start 启动server
func (s *server) Start() {
	// 获取listener
	serverTransport := s.transportMgr.Gen(s.opts.Protocol)
	ln, err := serverTransport.Listen(fmt.Sprintf("%v:%v", s.opts.IP, s.opts.Port))
	if err != nil {
		panic(err)
	}

	s.opts.Logger.Infof(context.TODO(), "server start at %s", fmt.Sprintf("%v:%v", s.opts.IP, s.opts.Port))

	// accept请求
	go s.run(ln)

	// 注册服务
	for name, _ := range s.serviceMap {
		s.opts.Discovery.Register(context.TODO(), &discov.Service{
			Name: name,
			Endpoints: []*discov.Endpoint{
				{
					ServiceName: name,
					IP:          s.opts.IP,
					Port:        s.opts.Port,
					Enable:      true,
				},
			},
		})
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Stop()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func (s *server) Stop() {
	// 服务取消注册
	for name, _ := range s.serviceMap {
		s.opts.Discovery.UnRegister(context.TODO(), &discov.Service{
			Name: name,
			Endpoints: []*discov.Endpoint{
				{
					ServiceName: name,
					IP:          s.opts.IP,
					Port:        s.opts.Port,
					Enable:      true,
				},
			},
		})
	}

	s.cancel()
}

func (s *server) run(ln net.Listener) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]
			s.opts.Logger.Errorf(context.TODO(), "err:%v\n.stack:%v", r, string(stack))
		}
	}()

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

			// 开启一个协程处理请求
			go s.process(accept)
		}

	}

}

// process logic
func (s *server) process(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]
			s.opts.Logger.Errorf(context.TODO(), "err:%v\n.stack:%v", r, string(stack))
		}
	}()

	//1.提取数据
	msg, err := s.extractMessage(conn)
	if err != nil {
		s.opts.Logger.Errorf(context.TODO(), "extractMessage error:%v", err.Error())
		return
	}

	//2.找到注册的服务和方法
	srv, ok := s.serviceMap[msg.ServiceName]
	if !ok {
		// todo 没有的service
	}

	methodHandler, ok := srv.methods[msg.ServiceMethod]
	if !ok {
		// todo 没有的方法
	}

	// 3.执行具体方法（包括注册的中间件）
	resp, err := methodHandler(context.TODO(), msg.Payload)

	// 4.对client发送返回数据
	var response codec.Response
	if err != nil {
		response.Msg = err.Error()
		response.Code = -1
	} else {
		response.Msg = ""
		response.Code = 0
	}
	response.Data = resp
	raw, _ := serialize.GetSerialize(s.opts.Serialize).Marshal(&response)
	if _, err = conn.Write(raw); err != nil {
		//s.opts.Logger.Errorf(context.TODO(), "extractMessage error:%v", err.Error())
	}

	return
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
