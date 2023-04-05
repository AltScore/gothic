package xgrpc

import (
	"context"
	"github.com/AltScore/auth-api/lib/auth"
	"github.com/AltScore/auth-api/lib/jwt"
	echojwt "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const JwtMetadataKey = "authorization"

func NewUserServerInterceptor(jwtManager jwt.Manager, next grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, found := metadata.FromIncomingContext(ctx)

		if found {
			if jwtStr, ok := md[JwtMetadataKey]; ok && len(jwtStr) > 0 {

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

// NewUserClientInterceptor is a gRPC client-side interceptor that adds the user credentials to the output request context.
func NewUserClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// Take user from context
		tokenObj := ctx.Value(auth.JwtKey)

		token, ok := tokenObj.(*echojwt.Token)

		if ok {
			// Add token to metadata
			md := metadata.Pairs(JwtMetadataKey, token.Raw)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		// call invoker
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
