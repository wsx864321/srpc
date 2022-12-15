package lb

import "github.com/wsx864321/srpc/discov"

type LoadBalance interface {
	Name() string
	Pick(*discov.Service) (*discov.Endpoint, error)
}

var lbMgr = map[string]LoadBalance{
	LoadBalanceRandom: NewRandom(),
}

// RegisterLB 注册loadBalance
func RegisterLB(lbName string, lb LoadBalance) {
	lbMgr[lbName] = lb
}

// GetLB 获取loadBalance方式
func GetLB(lbName string) LoadBalance {
	return lbMgr[lbName]
}
