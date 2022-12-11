package etcd

import (
	"context"
	"testing"
	"time"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/stretchr/testify/assert"
)

func TestNewETCDRegister(t *testing.T) {
	NewETCDRegister()
}

func TestRegister_Register(t *testing.T) {
	register := NewETCDRegister(WithEndpoints([]string{"127.0.0.1:2371"}))

	service := &discov.Service{
		Name: "test",
		Endpoints: []*discov.Endpoint{
			{
				ServiceName: "test",
				IP:          "127.0.0.1",
				Port:        9557,
				Enable:      true,
			},
		},
	}
	register.Register(context.TODO(), service)
	time.Sleep(2 * time.Second)
	registerService := register.GetService(context.TODO(), "test")

	assert.Equal(t, *service.Endpoints[0], *registerService.Endpoints[0])
}
