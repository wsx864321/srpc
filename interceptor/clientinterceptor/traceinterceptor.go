package clientinterceptor

import (
	"context"

	srpcerr "github.com/wsx864321/srpc/err"

	"github.com/wsx864321/srpc/interceptor"
	"github.com/wsx864321/srpc/metadata"
	strace "github.com/wsx864321/srpc/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ClientTraceInterceptor trace interceptor
func ClientTraceInterceptor() interceptor.ClientInterceptor {
	return func(ctx context.Context, method, target string, req, resp interface{}, h interceptor.Invoker) error {
		md := metadata.ExtractClientMetadata(ctx)

		tr := otel.GetTracerProvider().Tracer(strace.TraceName)
		name, attrs := strace.BuildSpan(method, target)
		ctx, span := tr.Start(ctx, name, trace.WithAttributes(attrs...), trace.WithSpanKind(trace.SpanKindClient))
		defer span.End()

		strace.Inject(ctx, otel.GetTextMapPropagator(), md)
		ctx = metadata.WithClientMetadata(ctx, md)

		err := h(ctx, req, resp)
		if err != nil {
			s, ok := srpcerr.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.GetMsg())
				span.SetAttributes(strace.StatusCodeAttr(s.GetCode()))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
			return err
		}

		span.SetAttributes(strace.StatusCodeAttr(0))
		return nil
	}
}
