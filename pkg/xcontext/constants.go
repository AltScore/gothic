package xcontext

const (
	// JwtCtxKey is the key used to store the JWT (raw) string in the context
	JwtCtxKey = "as-jwt-token"
	// UserCtxKey is the key used to store the user in the context
	UserCtxKey = "x-user"
	// ImpersonatingUserKeyName is the key used to store the impersonating (real)
	// user in the context
	ImpersonatingUserKeyName = "x-real-user"
	// TenantCtxKey is the key used to store the tenant in the context. If present
	// it will be used in case the user is not present.
	TenantCtxKey = "x-tenant"

	// DefaultTenant is the tenant to use if no tenant is found in the context
	DefaultTenant = "default"
)
