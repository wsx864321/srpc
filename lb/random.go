package lb

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/wsx864321/sweet_rpc/discov"
)

type Random struct {
}

func NewRandom() LoadBalance {
	return &Random{}
}

func (r *Random) Name() string {
	return "random"
}

func (r *Random) Pick(service *discov.Service) (*discov.Endpoint, error) {
	count := len(service.Endpoints)
	if count == 0 {
		return nil, errors.New("endpoint is empty")
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(count)))

	return service.Endpoints[n.Int64()], nil
}
