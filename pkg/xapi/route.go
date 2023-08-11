package xapi

import (
	"github.com/labstack/echo/v4"
)

// noPermissions is a convenience variable for routes that do not require any permissions
var noPermissions []string = nil

type Route struct {
	Method      string   // HTTP method
	Path        string   // path to the route with path parameters
	Permissions []string // user should have ALL of these permissions to access the route
	Hidden      bool     // if true, the route will not be listed in the documentation nor logs
	IsPublic    bool     // if true, the route will be accessible without authentication
	Handler     echo.HandlerFunc
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
