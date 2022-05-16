package etcd

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/wsx864321/sweet_rpc/discov"

	"github.com/stretchr/testify/assert"
)

func TestNewETCDRegister(t *testing.T) {
	register := NewETCDRegister()
	err := register.Init(context.TODO())

	assert.Nil(t, err)
}

func TestRegister_Register(t *testing.T) {
	register := NewETCDRegister()
	register.Init(context.TODO())

	service := &discov.Service{
		Name: "test",
		Endpoints: []*discov.Endpoint{
			&discov.Endpoint{
				ServerName: "test",
				IP:         "127.0.0.1",
				Port:       9557,
				Protocol:   1,
				Enable:     true,
			},
		},
	}
	register.Register(context.TODO(), service)
	time.Sleep(2 * time.Second)
	registerService := register.GetService(context.TODO(), "test")
	fmt.Println(*registerService.Endpoints[0])
}
