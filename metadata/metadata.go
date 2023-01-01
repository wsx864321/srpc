package metadata

import "context"

const (
	SRPCTimeout    = "srpc-timeout"
	SRPCPeerAddr   = "srpc-peer-addr"
	SPRCFullMethod = "srpc-full-method"
)

type clientMD struct{}
type serverMD struct{}

type ClientMetadata map[string]string

type ServerMetadata map[string]string

// ExtractClientMetadata extract ClientMetadata from context
func ExtractClientMetadata(ctx context.Context) ClientMetadata {
	if md, ok := ctx.Value(clientMD{}).(ClientMetadata); ok {
		return md
	}
	md := make(map[string]string)
	WithClientMetadata(ctx, md)
	return md
}

// WithClientMetadata creates a new context with the specified metadata
func WithClientMetadata(ctx context.Context, metadata map[string]string) context.Context {
	return context.WithValue(ctx, clientMD{}, ClientMetadata(metadata))
}

// ExtractServerMetadata extract ServerMetadata from context
func ExtractServerMetadata(ctx context.Context) ServerMetadata {
	if md, ok := ctx.Value(serverMD{}).(ServerMetadata); ok {
		return md
	}
	md := make(map[string]string)
	WithServerMetadata(ctx, md)
	return md
}

// WithServerMetadata creates a new context with the specified metadata
func WithServerMetadata(ctx context.Context, metadata map[string]string) context.Context {
	return context.WithValue(ctx, serverMD{}, ServerMetadata(metadata))
}
