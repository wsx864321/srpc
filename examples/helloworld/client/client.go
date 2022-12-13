package main

import (
	"encoding/json"
	"fmt"
	"github.com/wsx864321/srpc/codec"
	"net"
)

type HelloWorldReq struct {
	Name string `json:"name"`
}

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:9557")
	if err != nil {
		fmt.Println(err.Error(), "111")
		return
	}

	raw, _ := json.Marshal(&HelloWorldReq{
		Name: "wsx",
	})
	c := codec.NewCodec()
	req, err := c.Encode(1, 1, 11111, []byte("helloworld"), []byte("SayHello"), []byte("matedata"), raw)
	conn.Write(req)
	buf := make([]byte, 100000)
	conn.Read(buf)
	fmt.Println(string(buf))
}
