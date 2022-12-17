package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/discov/etcd"
)

func main() {
	req := &HelloWorldReq{
		Name: "wsx",
	}
	var resp HelloWorldResp
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	cli := client.NewClient(client.WithServiceName("helloworld"), client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	err := cli.Call(ctx, "SayHello", req, &resp)
	fmt.Println(resp, err)
}
