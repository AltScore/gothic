package xapi

import (
	"context"
	"github.com/labstack/echo/v4"
)

// apiContext is a wrapper around the echo.Context that implements the context.Context interface.
// This allows us to use the echo.Context as a context.Context.
// https://pkg.go.dev/context
// https://go.dev/blog/context-and-structs
// ignore lint
//nolint:containedctx // This is valid context
type apiContext struct {
	context.Context
	apiCtx echo.Context
}

func FromApi(ctx echo.Context) context.Context {
	return &apiContext{
		Context: ctx.Request().Context(),
		apiCtx:  ctx,
	}
}

func (a apiContext) Value(key any) any {
	if keyStr, ok := key.(string); ok {
		if val := a.apiCtx.Get(keyStr); val != nil {
			return val
		}
	}

	return a.Context.Value(key)
}
