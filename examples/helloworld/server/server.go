package main

import (
	"github.com/wsx864321/sweet_rpc/discov/etcd"
	"github.com/wsx864321/sweet_rpc/server"
)

func main() {
	s := server.NewServer(server.WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	s.RegisterService("helloworld", &HelloWorld{})
	s.Start()
}
