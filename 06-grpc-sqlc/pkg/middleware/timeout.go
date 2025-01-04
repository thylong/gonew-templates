package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// DefaultTimeoutUnaryInterceptor sets a default timeout for gRPC requests.
func DefaultTimeoutUnaryInterceptor(defaultTimeout time.Duration) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check if a deadline is already set
		if _, ok := ctx.Deadline(); !ok {
			// Set the default timeout
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, defaultTimeout)
			defer cancel()
		}

		// Call the handler
		return handler(ctx, req)
	}
}
