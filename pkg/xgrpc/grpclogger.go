package xgrpc

import (
	"context"
	"github.com/AltScore/gothic/pkg/xlogger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

// NewLoggerInterceptor is a gRPC server-side interceptor that logs requests
func NewLoggerInterceptor(logger xlogger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		result, err := handler(ctx, req)

		elapsed := time.Since(start)

		logger.Info("GRPC request", zap.String("method", info.FullMethod), zap.Int64("elapsed-ns", elapsed.Nanoseconds()), zap.Error(err))

		if err != nil {
			return result, convertError(err)
		}

		return result, err
	}
}
