package main

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/interceptor/clientinterceptor"
	strace "github.com/wsx864321/srpc/trace"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/pool"
)

func BenchmarkClient(b *testing.B) {
	strace.StartAgent(strace.WithServiceName("helloworld-client"))
	defer strace.StopAgent()

	req := &HelloWorldReq{
		Name: "wsx",
	}

	cli := client.NewClient(
		client.WithServiceName("helloworld"),
		client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))),
		client.WithPool(pool.NewPool(pool.WithInitialCap(10), pool.WithMaxCap(100))),
		client.WithReadTimeout(5*time.Second),
		client.WithWriteTimeout(5*time.Second),
		client.WithInterceptors([]interceptor.ClientInterceptor{clientinterceptor.ClientTraceInterceptor(), clientinterceptor.ClientTimeoutInterceptor()}...),
	)
	defer cli.Close()

	loopCount := 100000
	var count = 0
	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(loopCount)
	for i := 0; i < loopCount; i++ {
		go func() {
			defer wg.Done()
			var resp HelloWorldResp
			ctx, _ := context.WithTimeout(context.TODO(), 4*time.Second)
			err := cli.Call(ctx, "SayHello", req, &resp)
			if err != nil {
				count++
				fmt.Println(err)
			}
		}()

	}

	wg.Wait()
	fmt.Println(time.Now().Sub(now).Milliseconds(), "===", count)
}
