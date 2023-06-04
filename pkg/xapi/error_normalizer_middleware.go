package xapi

import (
	"errors"
	"net/http"

	"github.com/AltScore/gothic/pkg/xerrors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type HTTPStatusProvider interface {
	HTTPStatus() int
}

type fieldError struct {
	Tag             string
	Namespace       string
	StructNamespace string
	Field           string
	Value           interface{}
	Param           string
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

			var xerror xerrors.HttpError
			if errors.As(err, &xerror) {
				return &echo.HTTPError{
					Code:     xerror.HTTPStatus(),
					Message:  getMessageFromHttpXError(xerror, err),
					Internal: err,
				}
			}

			var echoError *echo.HTTPError
			if errors.As(err, &echoError) {
				return &echo.HTTPError{
					Code:     echoError.Code,
					Message:  getMessageFromHTTPError(echoError, echoError.Code),
					Internal: err,
				}
			}

			var statusProvider HTTPStatusProvider
			if errors.As(err, &statusProvider) {
				status := statusProvider.HTTPStatus()
				return &echo.HTTPError{
					Code:     status,
					Message:  getMessageFromError(err, status),
					Internal: err,
				}
			}

			var validationErr validator.FieldError
			if errors.As(err, &validationErr) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  messageFromValidationErr(validationErr),
					Internal: err,
				}
			}

			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				return &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  messageFromValidationErrors(validationErrors),
					Internal: err,
				}
			}

			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  getMessageFromError(err, http.StatusInternalServerError),
				Internal: err,
			}
		}
	}
}

func getMessageFromHTTPError(err *echo.HTTPError, code int) interface{} {
	return echo.Map{
		"error": echo.Map{
			"code":    xerrors.FromHttpStatus(code).(xerrors.HttpError).Code(),
			"message": err.Message,
		},
	}
}

func getMessageFromHttpXError(err xerrors.HttpError, err2 error) interface{} {
	return echo.Map{
		"error": echo.Map{
			"code":    err.Code(),
			"message": err2.Error(),
		},
	}
}

func messageFromValidationErrors(validationErrors validator.ValidationErrors) interface{} {
	fieldErrors := make([]*fieldError, 0, len(validationErrors))
	for _, err := range validationErrors {
		fieldErrors = append(fieldErrors, fieldErrorFromError(err))
	}

	return echo.Map{
		"error": echo.Map{
			"code":    "validation_error",
			"message": "Validation failed",
			"details": fieldErrors,
		},
	}
}

func messageFromValidationErr(err validator.FieldError) interface{} {
	return echo.Map{
		"error": echo.Map{
			"code":    err.Tag(),
			"message": err.Error(),
			"details": fieldErrorFromError(err),
		},
	}
}

func getMessageFromError(err error, status int) interface{} {
	return echo.Map{
		"error": echo.Map{
			"code":    xerrors.FromHttpStatus(status).(xerrors.HttpError).Code(),
			"message": err.Error(),
		},
	}
}

func fieldErrorFromError(fe validator.FieldError) *fieldError {
	return &fieldError{
		Tag:             fe.Tag(),
		Namespace:       fe.Namespace(),
		StructNamespace: fe.StructNamespace(),
		Field:           fe.Field(),
		Value:           fe.Value(),
		Param:           fe.Param(),
	}
}
