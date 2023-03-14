package xcontext

import (
	"context"
	"github.com/AltScore/gothic/pkg/xuser"
	"github.com/AltScore/lccb-api/pkg/xerrors"
)

var unauthorized = xerrors.NewUnauthorized("unauthorized")

const (
	UserCtxKey    = "x-user"
	TenantCtxKey  = "tenantID"
	DefaultTenant = "default"
)

// User returns the user from the context if it exists.
// If the user does not exist, returns an unauthorized error.
func User(ctx context.Context) (xuser.User, error) {
	user, found := _user(ctx)
	if !found {
		return nil, unauthorized
	}
	return user, nil
}

func _user(ctx context.Context) (xuser.User, bool) {
	user, found := ctx.Value(UserCtxKey).(xuser.User)
	return user, found
}

// Tenant returns the tenant from the context if it exists.
func Tenant(ctx context.Context) (tenant string, found bool) {
	user, found := _user(ctx)

	if found {
		return user.TenantID(), true
	}

	tenant, found = ctx.Value(TenantCtxKey).(string)
	return tenant, found
}

// TenantOrDefault returns the tenant from the context if it exists or the default value.
func TenantOrDefault(ctx context.Context) string {
	tenant, found := Tenant(ctx)
	if found {
		return tenant
	}
	return DefaultTenant
}
