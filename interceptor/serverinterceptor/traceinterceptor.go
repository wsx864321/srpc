package serverinterceptor

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

func ServerTraceInterceptor() interceptor.ServerInterceptor {
	return func(ctx context.Context, req interface{}, h interceptor.Handler) (interface{}, error) {
		md := metadata.ExtractServerMetadata(ctx)

		spanCtx := strace.Extract(ctx, otel.GetTextMapPropagator(), md)
		tr := otel.Tracer(strace.TraceName)
		name, attrs := strace.BuildSpan(md[metadata.SPRCFullMethod], md[metadata.SRPCPeerAddr])

		ctx, span := tr.Start(trace.ContextWithRemoteSpanContext(ctx, spanCtx), name, trace.WithSpanKind(trace.SpanKindServer), trace.WithAttributes(attrs...))
		defer span.End()

		resp, err := h(ctx, req)
		if err != nil {
			s, ok := srpcerr.FromError(err)
			if ok {
				span.SetStatus(codes.Error, s.GetMsg())
				span.SetAttributes(strace.StatusCodeAttr(s.GetCode()))
			} else {
				span.SetStatus(codes.Error, err.Error())
			}
			return nil, err
		}

		return resp, nil
	}
}
