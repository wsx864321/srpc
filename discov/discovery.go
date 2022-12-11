package discov

import "context"

type Discovery interface {
	Name() string
	Register(ctx context.Context, service *Service)
	UnRegister(ctx context.Context, service *Service)
	GetService(ctx context.Context, name string) *Service
}
