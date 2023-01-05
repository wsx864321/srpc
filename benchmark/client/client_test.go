package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/interceptor/clientinterceptor"
	"github.com/wsx864321/srpc/pool"
)

const cliCount = 100

func BenchmarkClient(b *testing.B) {
	req := &HelloWorldReq{
		Name: "wsx",
	}
	clis := make([]*client.Client, cliCount)
	for i := 0; i < cliCount; i++ {
		clis[i] = client.NewClient(
			client.WithServiceName("helloworld"),
			client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))),
			client.WithPool(pool.NewPool(pool.WithInitialCap(10), pool.WithMaxCap(100))),
			client.WithReadTimeout(5*time.Second),
			client.WithWriteTimeout(5*time.Second),
			client.WithInterceptors([]interceptor.ClientInterceptor{clientinterceptor.ClientTimeoutInterceptor()}...),
		)
	}

	defer func() {
		for _, c := range clis {
			c.Close()
		}
	}()

	reqCount := 1000000
	var count int64 = 0
	now := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(cliCount)
	for i := 0; i < cliCount; i++ {
		go func(cc int) {
			defer wg.Done()
			for j := 0; j < reqCount/cliCount; j++ {
				var resp HelloWorldResp
				ctx, _ := context.WithTimeout(context.TODO(), 4*time.Second)
				err := clis[cc].Call(ctx, "SayHello", req, &resp)
				if err != nil {
					atomic.AddInt64(&count, 1)
				}
			}

		}(i)

	}

	wg.Wait()
	fmt.Println(time.Now().Sub(now).Milliseconds(), "===", count)
}
