package xcontext

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/xerrors"
	"github.com/AltScore/gothic/v2/pkg/xuser"
)

var unauthorized = xerrors.NewUnauthorized("unauthorized")

// GetUser returns the user from the context if it exists.
// If the user does not exist, returns an unauthorized error.
func GetUser(ctx context.Context) (xuser.User, error) {
	user, found := getUser(ctx)
	if !found {
		return nil, unauthorized
	}
	return user, nil
}

func getUser(ctx context.Context) (xuser.User, bool) {
	user, isOk := ctx.Value(UserCtxKey).(xuser.User)
	return user, isOk
}

// GetTenant returns the tenant from the context if it exists.
func GetTenant(ctx context.Context) (tenant string, found bool) {
	user, found := getUser(ctx)

	if found {
		return user.Tenant(), true
	}

	tenant, found = ctx.Value(TenantCtxKey).(string)
	return tenant, found
}

// GetTenantOrDefault returns the tenant from the context if it exists or the default value.
func GetTenantOrDefault(ctx context.Context) string {
	tenant, found := GetTenant(ctx)
	if found {
		return tenant
	}
	return DefaultTenant
}

// MustGetUser returns the user from the context if it exists.
// If the user does not exist, it panics.
func MustGetUser(ctx context.Context) xuser.User {
	user, err := GetUser(ctx)
	if err != nil {
		panic(err)
	}
	return user
}

// WithUser returns a new context with the user set.
func WithUser(ctx context.Context, user xuser.User) context.Context {
	return context.WithValue(ctx, UserCtxKey, user)
}

// WithTenant returns a new context with the tenant set.
func WithTenant(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, TenantCtxKey, tenant)
}

// WithJwt returns a new context with the jwt set.
func WithJwt(ctx context.Context, jwt string) context.Context {
	return context.WithValue(ctx, JwtCtxKey, jwt)
}
