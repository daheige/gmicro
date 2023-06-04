package gmicro

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// CtxKey ctx key type.
type CtxKey string

const (
	// XRequestID request_id
	XRequestID CtxKey = "x-request-id"

	// GRPCClientIP grpc client_ip
	GRPCClientIP CtxKey = "client-ip"

	// RequestMethod request method
	RequestMethod CtxKey = "request_method"

	// RequestURI request uri
	RequestURI CtxKey = "request_uri"
)

// String returns string
func (c CtxKey) String() string {
	return string(c)
}

// GetIncomingMD returns metadata.MD from incoming ctx
// get request metadata
// this method is mainly used at the server end to get the relevant metadata data
func GetIncomingMD(ctx context.Context) metadata.MD {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md
}

// GetOutgoingMD returns metadata.MD from outgoing ctx
// Use this method when you pass ctx to a downstream service
func GetOutgoingMD(ctx context.Context) metadata.MD {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return metadata.MD{}
	}

	return md
}

// GetSliceFromMD returns []string from md
func GetSliceFromMD(md metadata.MD, key CtxKey) []string {
	return md.Get(key.String())
}

// GetSliceFromMD returns string from md
func GetStringFromMD(md metadata.MD, key CtxKey) string {
	values := md.Get(key.String())
	if len(values) > 0 {
		return values[0]
	}

	return ""
}

// SetCtxValue returns ctx when you set key/value into ctx
func SetCtxValue(ctx context.Context, key CtxKey, val interface{}) context.Context {
	return context.WithValue(ctx, key.String(), val)
}

// GetCtxValue returns ctx when you set key/value into ctx
func GetCtxValue(ctx context.Context, key CtxKey) interface{} {
	return ctx.Value(key)
}
