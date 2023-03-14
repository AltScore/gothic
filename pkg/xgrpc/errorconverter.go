package xgrpc

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type HTTPStatusProvider interface {
	HTTPStatus() int
}

// convertError converts errors to gRPC errors
func convertError(err error) error {
	var statusProvider HTTPStatusProvider
	if errors.As(err, &statusProvider) {
		return convertWithHttpStatus(err, statusProvider)
	}

	var validationErr validator.FieldError
	if errors.As(err, &validationErr) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return err
}

func convertWithHttpStatus(err error, provider HTTPStatusProvider) error {
	switch provider.HTTPStatus() {
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, err.Error())
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, err.Error())
	case http.StatusForbidden:
		return status.Error(codes.PermissionDenied, err.Error())
	case http.StatusNotFound:
		return status.Error(codes.NotFound, err.Error())
	case http.StatusConflict:
		return status.Error(codes.AlreadyExists, err.Error())
	case http.StatusUnprocessableEntity:
		return status.Error(codes.InvalidArgument, err.Error())
	case http.StatusTooManyRequests:
		return status.Error(codes.ResourceExhausted, err.Error())
	case http.StatusInternalServerError:
		return status.Error(codes.Internal, err.Error())
	case http.StatusServiceUnavailable:
		return status.Error(codes.Unavailable, err.Error())

	default:
		return status.Error(codes.Unknown, err.Error())
	}
}
