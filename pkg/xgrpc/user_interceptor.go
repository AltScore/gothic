package xgrpc

import (
	"context"
	"github.com/AltScore/auth-api/lib/auth"
	"github.com/AltScore/auth-api/lib/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func NewUserServerInterceptor(jwtManager jwt.Manager, next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, found := metadata.FromIncomingContext(ctx)

		if found {
			if jwtStr, ok := md["jwt"]; ok && len(jwtStr) > 0 {

				user, err := jwtManager.UserFromJwt(jwtStr[0])

				if err != nil {
					return nil, err
				}

				expandedCtx := context.WithValue(ctx, auth.UserCtxKey, user)
				return next(expandedCtx, req, info, handler)
			}
		}

		return next(ctx, req, info, handler)
	}
}
