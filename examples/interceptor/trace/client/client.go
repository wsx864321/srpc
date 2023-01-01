package main

import (
	"context"
	"fmt"
	"time"

	strace "github.com/wsx864321/srpc/trace"

	"github.com/wsx864321/srpc/interceptor/clientinterceptor"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/interceptor"
)

type HelloWorldReq struct {
	Name string `json:"name"`
}

type HelloWorldResp struct {
	Msg string `json:"msg"`
}

func main() {
	strace.StartAgent(strace.WithServiceName("helloworld-client"))
	defer strace.StopAgent()

	req := &HelloWorldReq{
		Name: "wsx",
	}
	var resp HelloWorldResp
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	cli := client.NewClient(
		client.WithServiceName("helloworld"),
		client.WithDiscovery(
			etcd.NewETCDRegister(
				etcd.WithEndpoints([]string{"127.0.0.1:2371"}),
			),
		),
		client.WithInterceptors([]interceptor.ClientInterceptor{clientinterceptor.ClientTraceInterceptor(), clientinterceptor.ClientTimeoutInterceptor()}...),
		client.WithReadTimeout(10*time.Second),
		client.WithWriteTimeout(10*time.Second),
	)
	err := cli.Call(ctx, "SayHello", req, &resp)
	fmt.Println(resp, err)
}
