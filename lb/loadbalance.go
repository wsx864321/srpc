package lb

import "github.com/wsx864321/sweet_rpc/discov"

type LoadBalanceType uint8


type LoadBalanceItf interface {
	Name() string
	Pick() (discov.Endpoint, error)
}
