package main

import (
	"github.com/wsx864321/srpc/codec/serialize"
	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/server"
)

func main() {
	s := server.NewServer(server.WithSerialize(serialize.SerializeTypeMsgpack), server.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	s.RegisterService("helloworld", &HelloWorld{})
	s.Start()
}
