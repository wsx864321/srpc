package main

import (
	_ "net/http/pprof"

	"github.com/wsx864321/srpc/codec/serialize"

	"github.com/wsx864321/srpc/discov/etcd"
	"github.com/wsx864321/srpc/server"
)

func main() {
	pprof()

	s := server.NewServer(
		server.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))),
		server.WithSerialize(serialize.SerializeTypeProto),
	)
	s.RegisterService("helloworld", &HelloWorld{})
	s.Start()
}
