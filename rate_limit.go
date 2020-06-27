package gmicro

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Limiter defines the interface to perform request rate limiting.
// If Limit function return true, the request will be rejected.
// Otherwise, the request will pass.
type Limiter interface {
	Limit() bool
}

// WithRateLimit rate limit
func WithRateLimit(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Limit() {
			return nil, status.Errorf(codes.ResourceExhausted,
				"%s is rejected by rate_limit,please retry later.", info.FullMethod)
		}

		return handler(ctx, req)
	}
}
