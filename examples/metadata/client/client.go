package main

import (
	"context"
	"fmt"
	"time"

	"github.com/wsx864321/srpc/metadata"

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
	var resp HelloWorldResp
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	ctx = metadata.WithClientMetadata(ctx, map[string]string{
		"token": "test",
	})
	cli := client.NewClient(client.WithServiceName("helloworld"), client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	err := cli.Call(ctx, "SayHello", req, &resp)
	fmt.Println(resp, err)
}
