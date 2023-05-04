package xapi

import (
	"github.com/AltScore/auth-api/lib/auth"
	"github.com/labstack/echo/v4"
)

// noPermissions is a convenience variable for routes that do not require any permissions
var noPermissions []string = nil

type PrivateHandlerFunc func(echo.Context, auth.User) error

type Route struct {
	Method      string   // HTTP method
	Path        string   // path to the route with path parameters
	Permissions []string // user should have ALL of these permissions to access the route
	Hidden      bool     // if true, the route will not be listed in the documentation nor logs
	IsPublic    bool     // if true, the route will be accessible without authentication
	Handler     echo.HandlerFunc
}

// Private creates a route that requires authentication and has the given permissions
func Private(method string, path string, handler PrivateHandlerFunc, permissions ...string) Route {
	return Route{
		Method:      method,
		Path:        path,
		Permissions: permissions,
		IsPublic:    false,
		Handler:     adaptPrivateHandler(handler),
	}
}

func adaptPrivateHandler(handler PrivateHandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, found := auth.UserFromRequest(c)
		if !found {
			return echo.NewHTTPError(401, "user not found")
		}
		return handler(c, user)
	}
}

// Public creates a route that does not require authentication
func Public(method string, path string, handler echo.HandlerFunc) Route {

	return Route{
		Method:      method,
		Path:        path,
		IsPublic:    true,
		Permissions: noPermissions,
		Handler:     handler,
	}
}

// Hidden creates a route that does not require authentication and is hidden from the documentation, logs, and metrics
func Hidden(method string, path string, handler echo.HandlerFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		IsPublic:    true,
		Permissions: noPermissions,
		Hidden:      true,
		Handler:     handler,
	}
}
