package serverinterceptor

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/metadata"
)

const (
	sprcTimeout = "sprc-timeout"
)

// ServerTimeoutInterceptor 服务端超时控制，与client进行配合做到级联超时控制
func ServerTimeoutInterceptor() interceptor.ServerInterceptor {
	return func(ctx context.Context, req interface{}, h interceptor.Handler) (interface{}, error) {
		md := metadata.ExtractServerMetadata(ctx)
		if val, ok := md[sprcTimeout]; ok {
			dur, _ := strconv.ParseInt(val, 10, 64)
			ctx, _ = context.WithTimeout(ctx, time.Nanosecond*time.Duration(dur))
		}

		var (
			finish = make(chan struct{}, 1)
			resp   interface{}
			err    error
		)

		go func() {
			resp, err = h(ctx, req)

			finish <- struct{}{}
		}()

		// 执行业务逻辑后操作
		select {
		case <-finish:
			return resp, err
		case <-ctx.Done():
			return nil, errors.New("server timeout")
		}
	}
}
