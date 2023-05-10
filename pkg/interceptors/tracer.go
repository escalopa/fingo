package interceptors

import (
	"context"

	"github.com/escalopa/fingo/pkg/contextutils"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func TracingUnaryInterceptor() func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		_, err := contextutils.GetTracerID(ctx)
		// Check if tracer is already registered
		// If not, create one and register it in context
		if err != nil {
			newID := uuid.New().String()
			ctx = contextutils.SetTracerID(ctx, newID)
		}
		return handler(ctx, req)
	}
}
