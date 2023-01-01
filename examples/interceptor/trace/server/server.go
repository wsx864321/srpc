package main

import (
	"github.com/wsx864321/srpc/codec/serialize"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/interceptor/serverinterceptor"
	"github.com/wsx864321/srpc/server"
	strace "github.com/wsx864321/srpc/trace"
)

func main() {
	strace.StartAgent(strace.WithServiceName("helloworld-server"))
	defer strace.StopAgent()

	s := server.NewServer(
		server.WithSerialize(serialize.SerializeTypeMsgpack),
		server.WithDiscovery(
			etcd.NewETCDRegister(
				etcd.WithEndpoints([]string{"127.0.0.1:2371"}),
			),
		),
	)
	s.RegisterService("helloworld", &HelloWorld{}, []interceptor.ServerInterceptor{serverinterceptor.ServerTraceInterceptor()}...)
	s.Start()
}
