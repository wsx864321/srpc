package clientinterceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/metadata"
)

// ClientTimeoutInterceptor 客户端级联超时控制
func ClientTimeoutInterceptor() interceptor.ClientInterceptor {
	return func(ctx context.Context, method, target string, req, resp interface{}, h interceptor.Invoker) error {
		if deadline, ok := ctx.Deadline(); ok {
			md := metadata.ExtractClientMetadata(ctx)
			md[metadata.SRPCTimeout] = fmt.Sprintf("%d", time.Until(deadline).Nanoseconds())
			ctx = metadata.WithClientMetadata(ctx, md)
		}

		return h(ctx, req, resp)
	}
}
