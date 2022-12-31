package clientinterceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/metadata"
)

const (
	sprcTimeout = "sprc-timeout"
)

// ClientTimeoutInterceptor 客户端级联超市控制
func ClientTimeoutInterceptor() interceptor.ClientInterceptor {
	return func(ctx context.Context, req, resp interface{}, h interceptor.Invoker) error {
		if deadline, ok := ctx.Deadline(); ok {
			md := metadata.ExtractClientMetadata(ctx)
			md[sprcTimeout] = fmt.Sprintf("%d", time.Until(deadline).Nanoseconds())
			ctx = metadata.WithClientMetadata(ctx, md)
		}

		return h(ctx, req, resp)
	}
}
