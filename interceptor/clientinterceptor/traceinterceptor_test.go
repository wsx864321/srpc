package clientinterceptor

import (
	"context"
	"testing"

	strace "github.com/wsx864321/srpc/trace"
)

func TestClientTraceInterceptor(t *testing.T) {
	strace.StartAgent()
	defer strace.StopAgent()

	ClientTraceInterceptor()(context.Background(), "srpc-test/helloworld", "127.0.0.1:7777", "", "", func(ctx context.Context, req, resp interface{}) error {
		return nil
	})
}
