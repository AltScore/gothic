package xapi

import (
	"github.com/AltScore/gothic/v2/pkg/xcontext"
	"net/http"
	"net/url"

	"github.com/AltScore/gothic/v2/pkg/ids"
	"github.com/AltScore/gothic/v2/pkg/xerrors"
	"github.com/AltScore/gothic/v2/pkg/xuser"
	"github.com/labstack/echo/v4"
)

var unauthorized = xerrors.NewUnauthorized("unauthorized")

// UserFromContext returns the authenticated user from the context.
// If the user does not exist, or it is an invalid struct, found will be false.
func UserFromContext(c echo.Context) (user xuser.User, found bool) {
	ctx := c.Request().Context()

	value := ctx.Value(xcontext.UserCtxKey)

	if value == nil {
		return nil, false
	}

	if apiUser, ok := value.(xuser.User); ok {
		return apiUser, true
	}

	return nil, false
}

// TenantFromContext returns the current user tenant from the API context.
func TenantFromContext(c echo.Context) (string, bool) {
	user, found := UserFromContext(c)

	if found {
		return user.Tenant(), true
	}

	return "", false
}

func TenantOrDefault(c echo.Context) string {
	tenant, found := TenantFromContext(c)

	if found {
		return tenant
	}

	return xcontext.DefaultTenant
}

// UserFromContextOrError returns the authenticated user from the context.
// If the user does not exist, or it is an invalid struct, returns an error.
func UserFromContextOrError(c echo.Context) (user xuser.User, err error) {
	value := c.Get(xcontext.UserCtxKey)

	if value == nil {
		return nil, unauthorized
	}

	if apiUser, ok := value.(xuser.User); ok {
		return apiUser, nil
	}

	return nil, unauthorized
}

// RealUserFromContext returns the real authenticated user.
// If user was impersonated, it returns the impersonating user. If not, the same user in context is returned.
func RealUserFromContext(c echo.Context) (xuser.User, bool) {
	value := c.Get(xcontext.ImpersonatingUserKeyName)

	if value == nil {
		return UserFromContext(c)
	}

	if user, ok := value.(xuser.User); ok {
		return user, true
	}

	return UserFromContext(c)
}

func ParseParamID(c echo.Context, name string) (ids.Id, error) {
	idStr, err := url.PathUnescape(c.Param(name))

	if err != nil {
		return ids.Empty(), err
	} else if idStr == "" {
		return ids.Empty(), echo.NewHTTPError(http.StatusBadRequest, name+" is required")
	}

	return ids.Parse(idStr)
}
