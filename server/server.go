package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"syscall"
	"time"

	"github.com/wsx864321/srpc/metadata"

	"github.com/gogo/protobuf/proto"
	"github.com/wsx864321/srpc/codec"
	"github.com/wsx864321/srpc/codec/serialize"
	"github.com/wsx864321/srpc/discov"
	srpcerr "github.com/wsx864321/srpc/err"
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/transport"
	"github.com/wsx864321/srpc/util"
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
	codec        *codec.Codec
	ctx          context.Context    // Each service is managed in one context
	cancel       context.CancelFunc // controller of context
	interceptors []interceptor.ServerInterceptor
	serviceMap   map[string]*service
}

type service struct {
	name    string
	methods map[string]handler
}

// NewServer 生成一个server
func NewServer(opts ...Option) Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		opts:         NewOptions(opts...),
		codec:        codec.NewCodec(),
		ctx:          ctx,
		cancel:       cancel,
		interceptors: make([]interceptor.ServerInterceptor, 0),
		serviceMap:   make(map[string]*service),
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
		if err := s.checkMethod(method.Name, method.Type); err != nil {
			panic(err)
		}

		methodHandler := func(ctx context.Context, body []byte) (interface{}, error) {
			req := reflect.New(method.Type.In(2).Elem()).Interface()
			if err := serialize.GetSerialize(s.opts.Serialize).Unmarshal(body, req); err != nil {
				return nil, err
			}

			h := func(ctx context.Context, req interface{}) (interface{}, error) {
				resp := method.Func.Call([]reflect.Value{svrValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
				errInterface := resp[1].Interface()
				if errInterface != nil {
					return resp[0].Interface(), errInterface.(error)
				}
				return resp[0].Interface(), nil

			}

			return interceptor.ServerIntercept(ctx, req, s.interceptors, h)
		}
		methods[method.Name] = methodHandler
	}

	s.serviceMap[serName] = &service{
		serName,
		methods,
	}
}

// checkMethod
func (s *server) checkMethod(methodName string, methodType reflect.Type) error {
	if methodType.NumIn() != 3 {
		return fmt.Errorf("method %s invalid, the number of params != 2", methodName)
	}

	if methodType.NumOut() != 2 {
		return fmt.Errorf("method %s invalid, the number of params != 2", methodName)
	}

	var ctx *context.Context
	if !methodType.In(1).Implements(reflect.TypeOf(ctx).Elem()) {
		return fmt.Errorf("method %s invalid, first param is not context", methodName)
	}

	if s.opts.Serialize == serialize.SerializeTypeProto {
		var p *proto.Message
		if !methodType.In(2).Implements(reflect.TypeOf(p).Elem()) {
			return fmt.Errorf("method %s invalid, second param is not proto.Message", methodName)
		}
	} else {
		if methodType.In(2).Kind() != reflect.Ptr {
			return fmt.Errorf("method %s invalid, second param is not ptr", methodName)
		}
	}

	if s.opts.Serialize == serialize.SerializeTypeProto {
		var p *proto.Message
		if !methodType.Out(0).Implements(reflect.TypeOf(p).Elem()) {
			return fmt.Errorf("method %s invalid, second param is not proto.Message", methodName)
		}
	} else {
		if methodType.Out(0).Kind() != reflect.Ptr {
			return fmt.Errorf("method %s invalid, first reply type is not a pointer", methodName)
		}
	}

	var err *error
	if !methodType.Out(1).Implements(reflect.TypeOf(err).Elem()) {
		return fmt.Errorf("method %s invalid, second reply type is not error", methodName)
	}

	return nil
}

// Start 启动server
func (s *server) Start() {
	// 获取listener
	serverTransport := transport.GetTransport(s.opts.Network)
	ln, err := serverTransport.Listen(fmt.Sprintf("%v:%v", s.opts.IP, s.opts.Port))
	if err != nil {
		panic(err)
	}

	s.opts.Logger.Infof(context.TODO(), "%s server start at %s", s.opts.Network, fmt.Sprintf("%v:%v", s.opts.IP, s.opts.Port))

	// accept请求
	go s.run(ln)

	// 等待100ms，让服务先接收请求
	time.Sleep(100 * time.Millisecond)

	// 注册服务
	for name := range s.serviceMap {
		s.opts.Discovery.Register(context.TODO(), &discov.Service{
			Name: name,
			Endpoints: []*discov.Endpoint{
				{
					ServiceName: name,
					IP:          s.opts.IP,
					Port:        s.opts.Port,
					Network:     string(s.opts.Network),
					Serialize:   string(s.opts.Serialize),
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
	for name := range s.serviceMap {
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
			go s.handleConn(s.ctx, accept)
		}

	}

}

// handleConn logic
func (s *server) handleConn(pCtx context.Context, conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]
			s.opts.Logger.Errorf(context.TODO(), "\nerr:%v\nstack:%v", r, string(stack))
		}
	}()

	for {
		select {
		case <-pCtx.Done():
			return
		default:

		}

		if err := s.process(conn); err != nil {
			return
		}
	}

}

func (s *server) process(conn net.Conn) error {
	//1.提取数据
	msg, err := s.extractMessage(conn)
	// 这里为什么不用io.EOF 可以见：https://blog.csdn.net/aixinaxc/article/details/89282364
	if err != nil {
		s.opts.Logger.Errorf(context.TODO(), "extractMessage error:%v", err.Error())
		return err
	}

	// 2.提取metadata
	var metaData metadata.Metadata
	if len(msg.MetaData) != 0 {
		if err = serialize.GetSerialize(s.opts.Serialize).Unmarshal(msg.MetaData, &metaData); err != nil {
			s.wireErr(context.TODO(), conn, srpcerr.NewError(srpcerr.UnKnowErr, err.Error()))
			return nil
		}
	}
	ctx := metadata.WithServerMetadata(context.Background(), metaData.Data)

	//3.找到注册的服务和方法
	srv, ok := s.serviceMap[msg.ServiceName]
	if !ok {
		s.wireErr(ctx, conn, srpcerr.ServiceNotExistErr)
		return nil
	}

	methodHandler, ok := srv.methods[msg.ServiceMethod]
	if !ok {
		s.wireErr(ctx, conn, srpcerr.MethodNotExistErr)
		return nil
	}

	// 3.执行具体方法（包括注册的中间件）
	if s.opts.Timeout > 0 {
		// 超时控制统一在timeoutInterceptor中间件中执行
		ctx, _ = context.WithTimeout(ctx, s.opts.Timeout)
	}
	resp, err := methodHandler(ctx, msg.Payload)

	// 4.对client发送返回数据
	if err != nil {
		if err, ok := err.(*srpcerr.Error); ok {
			s.wireErr(ctx, conn, err)
			return nil
		}

		s.wireErr(ctx, conn, srpcerr.NewError(srpcerr.UnKnowErr, err.Error()))
		return nil
	}

	raw, _ := serialize.GetSerialize(s.opts.Serialize).Marshal(resp)
	raw, _ = serialize.GetSerialize(s.opts.Serialize).Marshal(srpcerr.OkErr.WithData(raw))
	if s.opts.WriteTimeout > 0 {
		if err = conn.SetWriteDeadline(time.Now().Add(s.opts.WriteTimeout)); err != nil {
			resp, _ := serialize.GetSerialize(s.opts.Serialize).Marshal(err)
			raw, err = s.codec.Encode(codec.GeneralMsgType, codec.CompressTypeNot, uint64(time.Now().Unix()), []byte(""), []byte(""), []byte(""), resp)
			if err != nil {
				s.opts.Logger.Errorf(ctx, "write error:%v", err.Error())
				return nil
			}
			if err = util.Write(conn, raw); err != nil {
				s.opts.Logger.Errorf(ctx, "write error:%v", err.Error())
				return nil
			}
		}
	}

	raw, err = s.codec.Encode(codec.GeneralMsgType, codec.CompressTypeNot, uint64(time.Now().Unix()), []byte(""), []byte(""), []byte(""), raw)
	if err = util.Write(conn, raw); err != nil {
		s.opts.Logger.Errorf(ctx, "write error:%v", err.Error())
	}

	return nil
}

// extractMessage 提取message内容
func (s *server) extractMessage(conn net.Conn) (*codec.Message, error) {
	// 1.读取头部内容
	headerData := make([]byte, s.codec.GetHeaderLength())
	if err := util.ReadFixData(conn, headerData); err != nil {
		return nil, err
	}

	header, err := s.codec.DecodeHeader(headerData)
	if err != nil {
		return nil, err
	}

	// 2.读取message内容
	body := make([]byte, s.codec.GetBodyLength(header))
	if err = util.ReadFixData(conn, body); err != nil {
		return nil, err
	}

	return s.codec.DecodeBody(header, body)
}

func (s *server) wireErr(ctx context.Context, conn net.Conn, err *srpcerr.Error) {
	if s.opts.WriteTimeout > 0 {
		if e := conn.SetWriteDeadline(time.Now().Add(s.opts.WriteTimeout)); e != nil {
			s.opts.Logger.Errorf(context.TODO(), "SetWriteDeadline error:%v", err.Error())
			resp, _ := serialize.GetSerialize(s.opts.Serialize).Marshal(err)
			raw, e := s.codec.Encode(codec.GeneralMsgType, codec.CompressTypeNot, uint64(time.Now().Unix()), []byte(""), []byte(""), []byte(""), resp)
			if e != nil {
				s.opts.Logger.Errorf(ctx, "write error:%v", e.Error())
				return
			}

			if e = util.Write(conn, raw); e != nil {
				s.opts.Logger.Errorf(ctx, "write error:%v", e.Error())
			}
			return
		}
	}

	resp, _ := serialize.GetSerialize(s.opts.Serialize).Marshal(err)
	raw, e := s.codec.Encode(codec.GeneralMsgType, codec.CompressTypeNot, uint64(time.Now().Unix()), []byte(""), []byte(""), []byte(""), resp)
	if e != nil {
		s.opts.Logger.Errorf(ctx, "write error:%v", e.Error())
		return
	}

	if e := util.Write(conn, raw); e != nil {
		s.opts.Logger.Errorf(ctx, "write error:%v", e.Error())
	}

	return
}
