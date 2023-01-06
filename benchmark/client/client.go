package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/interceptor/clientinterceptor"

	"github.com/wsx864321/srpc/pool"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
)

const cliCount = 100

var concurrency = flag.Int64("concurrency", 100, "concurrency")
var total = flag.Int64("total", 1000000, "total requests")

type Counter struct {
	Succ        int64 // 成功量
	Fail        int64 // 失败量
	Total       int64 // 总量
	Concurrency int64 // 并发量
	Cost        int64 // 总耗时 ms
}

func main() {
	flag.Parse()
	benchmark(*total, *concurrency)
}

func benchmark(total int64, concurrency int64) {
	perClientReqs := int(total / concurrency)

	counter := &Counter{
		Total:       int64(perClientReqs) * concurrency,
		Concurrency: concurrency,
		Fail:        0,
	}
	req := &HelloWorldReq{
		Name: "wsx",
	}
	clis := make([]*client.Client, concurrency)
	for i := 0; i < int(concurrency); i++ {
		clis[i] = client.NewClient(
			client.WithServiceName("helloworld"),
			client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))),
			client.WithPool(pool.NewPool(pool.WithInitialCap(10), pool.WithMaxCap(200))),
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

	startTime := time.Now().UnixNano()
	wg := sync.WaitGroup{}
	wg.Add(int(concurrency))
	for i := 0; i < int(concurrency); i++ {
		go func(cc int) {
			defer wg.Done()
			for j := 0; j < perClientReqs; j++ {
				var resp HelloWorldResp
				ctx, _ := context.WithTimeout(context.TODO(), 40*time.Second)
				err := clis[cc].Call(ctx, "SayHello", req, &resp)
				if err != nil {
					atomic.AddInt64(&counter.Fail, 1)
				}
			}

		}(i)

	}

	wg.Wait()
	counter.Succ = total - counter.Fail
	counter.Cost = (time.Now().UnixNano() - startTime) / 1000000

	fmt.Printf("took %d ms for %d requests \n", counter.Cost, counter.Total)
	fmt.Printf("sent     requests      : %d\n", counter.Total)
	fmt.Printf("received requests      : %d\n", atomic.LoadInt64(&counter.Succ)+atomic.LoadInt64(&counter.Fail))
	fmt.Printf("received requests succ : %d\n", atomic.LoadInt64(&counter.Succ))
	fmt.Printf("received requests fail : %d\n", atomic.LoadInt64(&counter.Fail))
	fmt.Printf("throughput  (TPS)      : %d\n", total*1000/counter.Cost)
}
