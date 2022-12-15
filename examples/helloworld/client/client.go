package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wsx864321/srpc/client"
	"github.com/wsx864321/srpc/codec"
	"github.com/wsx864321/srpc/discov/etcd"
	"time"
)

type HelloWorldReq struct {
	Name string `json:"name"`
}

func main() {
	//conn, err := net.Dial("tcp", "0.0.0.0:9557")
	//if err != nil {
	//	fmt.Println(err.Error(), "111")
	//	return
	//}
	//
	raw, _ := json.Marshal(&HelloWorldReq{
		Name: "wsx",
	})
	c := codec.NewCodec()
	req, _ := c.Encode(1, 1, 11111, []byte("helloworld"), []byte("SayHello"), []byte("matedata"), raw)
	//conn.Write(req)
	resp := make([]byte, 100000)
	//conn.Read(resp)
	//fmt.Println(string(resp))

	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	cli := client.NewClient(client.WithServiceName("helloworld"), client.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	err := cli.Call(ctx, "SayHello", req, resp)
	fmt.Println(string(resp), err)
}
