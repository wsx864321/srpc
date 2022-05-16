package discov

import "context"

type Discovery interface {
	Init(ctx context.Context) error
	Name() string
	Register(ctx context.Context, service *Service)
	UnRegister(ctx context.Context, service *discov.Service)
	GetService(ctx context.Context, name string) *Service
}
