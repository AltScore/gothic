package xerrors

import (
	"errors"
	"fmt"
	"net/http"
)

type HttpError struct {
	code       string
	msg        string
	httpStatus int
}

func New(code string, msg string, httpStatus int) error {
	return HttpError{
		code:       code,
		msg:        msg,
		httpStatus: httpStatus,
	}
}

func (e HttpError) Code() string {
	return e.code
}

func (e HttpError) Error() string {
	return e.msg
}

func (e HttpError) HTTPStatus() int {
	return e.httpStatus
}

func (e HttpError) Unwrap() error {
	return nil
}

var (
	ErrNotFound         = New("not-found", "not found for", http.StatusNotFound)
	ErrDuplicate        = New("duplicate", "duplicate", http.StatusConflict)
	ErrFoundMany        = New("found-many", "found many but one expected", http.StatusConflict)
	ErrTypeAssertion    = New("invalid-type", "type assertion failed", http.StatusInternalServerError)
	ErrUnknown          = New("unknown", "unknown error found", http.StatusInternalServerError)
	ErrInvalidArgument  = New("invalid-argument", "invalid argument", http.StatusBadRequest)
	ErrInvalidState     = New("invalid-state", "invalid state", http.StatusPreconditionFailed)
	ErrClientCanceled   = New("cancelled", "client cancelled", 460)
	ErrTimeout          = New("timeout", "timeout", http.StatusGatewayTimeout)
	ErrGateway          = New("gateway", "gateway", http.StatusBadGateway)
	ErrUnauthorized     = New("unauthorized", "user did not provided credentials", http.StatusUnauthorized)  // Not authenticated,
	ErrForbidden        = New("forbidden", "user is not allowed to perform operation", http.StatusForbidden) // Not enough permissions
	ErrInvalidEventType = New("invalid-event-type", "invalid event type", http.StatusInternalServerError)
	ErrConditionNotMet  = New("condition-not-met", "condition not met", http.StatusPreconditionFailed)
)

func FromHttpStatus(status int) error {
	switch status {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return ErrInvalidArgument
	case http.StatusConflict:
		return ErrInvalidState
	case http.StatusPreconditionFailed:
		return ErrConditionNotMet
	case http.StatusBadGateway:
		return ErrGateway
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusGatewayTimeout:
		return ErrTimeout
	default:
		return ErrUnknown
	}
}

func NewUnknownError(entity string, details string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s: %s", ErrUnknown, entity, fmt.Sprintf(keyFmt, args...), details)
}

func NewInvalidArgumentError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidArgument, entity, fmt.Sprintf(keyFmt, args...))
}

func NewNotFoundError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrNotFound, entity, fmt.Sprintf(keyFmt, args...))
}

func NewDuplicateError(entity string, details string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s: %s", ErrDuplicate, entity, fmt.Sprintf(keyFmt, args...), details)
}

func NewFoundManyError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrFoundMany, entity, fmt.Sprintf(keyFmt, args...))
}

func NewConditionNotMetError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrConditionNotMet, entity, fmt.Sprintf(keyFmt, args...))
}

func NewTypeAssertionError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrTypeAssertion, entity, fmt.Sprintf(keyFmt, args...))
}

func NewTimeoutError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrTimeout, entity, fmt.Sprintf(keyFmt, args...))
}

func NewGatewayError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrGateway, entity, fmt.Sprintf(keyFmt, args...))
}

func NewCancellationError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrClientCanceled, entity, fmt.Sprintf(keyFmt, args...))
}

func NewInvalidStateError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidState, entity, fmt.Sprintf(keyFmt, args...))
}

func NewInvalidEventTypeError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidEventType, entity, fmt.Sprintf(keyFmt, args...))
}

func NewForbiddenError(entity string, keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s: %s", ErrForbidden, entity, fmt.Sprintf(keyFmt, args...))
}

func NewUnauthorized(keyFmt string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrUnauthorized, fmt.Sprintf(keyFmt, args...))
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, ErrNotFound)
}
