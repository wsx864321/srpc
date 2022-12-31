package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/wsx864321/srpc/pool"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
)

type HelloWorldReq struct {
	Name string `json:"name"`
}

type HelloWorldResp struct {
	Msg string `json:"msg"`
}

func main() {
	req := &HelloWorldReq{
		Name: "wsx",
	}

	cli := client.NewClient(
		client.WithServiceName("helloworld"),
		client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))),
		client.WithPool(pool.NewPool(pool.WithInitialCap(10), pool.WithMaxCap(1000))),
		client.WithReadTimeout(3*time.Second),
		client.WithWriteTimeout(3*time.Second),
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
			ctx, _ := context.WithTimeout(context.TODO(), 20*time.Second)
			err := cli.Call(ctx, "SayHello", req, &resp)
			if err != nil {
				count++
				//fmt.Println(err)
			}
		}()

	}

	wg.Wait()
	fmt.Println(time.Now().Sub(now).Milliseconds(), "===", count)

}
