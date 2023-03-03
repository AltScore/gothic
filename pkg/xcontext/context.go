package xcontext

import (
	"context"
	"github.com/AltScore/gothic/pkg/xuser"
)

const (
	UserCtxKey   = "x-user"
	TenantCtxKey = "tenantID"
)

// User returns the user from the context if it exists.
// If the user does not exist, found will be false.
func User(ctx context.Context) (user xuser.User, found bool) {
	user, found = ctx.Value(UserCtxKey).(xuser.User)
	return user, found
}

// Tenant returns the tenant from the context if it exists.
func Tenant(ctx context.Context) (tenant string, found bool) {
	tenant, found = ctx.Value(TenantCtxKey).(string)
	return tenant, found
}

// TenantOrDefault returns the tenant from the context if it exists or the given default value.
func TenantOrDefault(ctx context.Context, defaultTenant string) string {
	tenant, found := Tenant(ctx)
	if found {
		return tenant
	}
	return defaultTenant
}
