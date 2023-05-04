package xapi

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HTTPStatusProvider interface {
	HTTPStatus() int
}

// ErrorNormalizerMiddleware is a middleware that convert general errors to HTTPError.
// This middleware should be after the logger middleware.
// The conversion should be performed by the CustomErrorHandler but the logger middleware
// will not log the real status code if it is not a HTTPError.
func ErrorNormalizerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			var statusProvider HTTPStatusProvider
			if errors.As(err, &statusProvider) {
				return &echo.HTTPError{
					Code:     statusProvider.HTTPStatus(),
					Message:  err.Error(),
					Internal: err,
				}
			}

			var validationErr validator.FieldError
			if errors.As(err, &validationErr) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  err.Error(),
					Internal: err,
				}
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  err.Error(),
					Internal: err,
				}
			}

			return err
		}
	}
}
