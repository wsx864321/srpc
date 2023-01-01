package trace

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/trace"
)

const (
	TraceName = "srpc-trace"
)

type metadataSupplier struct {
	metadata map[string]string
}

func (s *metadataSupplier) Get(key string) string {
	return s.metadata[key]
}

func (s *metadataSupplier) Set(key, value string) {
	s.metadata[key] = value
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(s.metadata))
	for key := range s.metadata {
		out = append(out, key)
	}

	return out
}

// Inject set cross-cutting concerns from the Context into the metadata.
func Inject(ctx context.Context, p propagation.TextMapPropagator, m map[string]string) {
	p.Inject(ctx, &metadataSupplier{
		metadata: m,
	})
}

// Extract reads cross-cutting concerns from the metadata into a Context.
func Extract(ctx context.Context, p propagation.TextMapPropagator, m map[string]string) sdktrace.SpanContext {
	ctx = p.Extract(ctx, &metadataSupplier{
		metadata: m,
	})

	return sdktrace.SpanContextFromContext(ctx)
}
