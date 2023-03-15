package xapi

import (
	"github.com/AltScore/gothic/pkg/ids"
	"github.com/AltScore/gothic/pkg/xerrors"
	"github.com/AltScore/gothic/pkg/xuser"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

const UserCtxKey = "x-user"

var unauthorized = xerrors.NewUnauthorized("unauthorized")

const ImpersonatingUserKeyName = "ImpersonatingUser"

// UserFromContext returns the authenticated user from the context.
// If the user does not exist, or it is an invalid struct, found will be false.
func UserFromContext(c echo.Context) (user xuser.User, found bool) {
	value := c.Get(UserCtxKey)

	if value == nil {
		return nil, false
	}

	if apiUser, ok := value.(xuser.User); ok {
		return apiUser, true
	}

	return nil, false
}

// UserFromContextOrError returns the authenticated user from the context.
// If the user does not exist, or it is an invalid struct, returns an error.
func UserFromContextOrError(c echo.Context) (user xuser.User, err error) {
	value := c.Get(UserCtxKey)

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
	value := c.Get(ImpersonatingUserKeyName)

	if value == nil {
		return UserFromContext(c)
	}

	if user, ok := value.(xuser.User); ok {
		return user, true
	}

	return UserFromContext(c)
}

func ParseParamId(c echo.Context, name string) (ids.ID, error) {
	idStr, err := url.PathUnescape(c.Param(name))

	if err != nil {
		return ids.Empty(), err
	} else if idStr == "" {
		return ids.Empty(), echo.NewHTTPError(http.StatusBadRequest, name+" is required")
	}

	return ids.ParseId(idStr)
}
