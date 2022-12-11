package lb

import "github.com/wsx864321/sweet_rpc/discov"

type LoadBalanceType uint8

type LoadBalance interface {
	Name() string
	Pick(*discov.Service) (*discov.Endpoint, error)
}
