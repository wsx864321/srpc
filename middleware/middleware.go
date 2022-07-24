package middleware

import "context"

type MiddlewareFunc func(ctx context.Context, req interface{}) (interface{}, error)
