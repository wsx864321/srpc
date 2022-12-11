package server

import (
	"github.com/wsx864321/srpc/discov/etcd"
	"testing"
)

func TestNewServer(t *testing.T) {
	s := NewServer(WithDiscovery(etcd.NewETCDRegister(etcd.WithEndpoints([]string{"127.0.0.1:2371"}))))
	s.Start()
}
