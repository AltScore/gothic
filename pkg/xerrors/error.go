package xerrors

import (
	"errors"
	"fmt"
)

const (
	NotFoundReason         = "not found for"
	DuplicateReason        = "duplicate"
	FoundManyReason        = "found many but one expected"
	TypeAssertionReason    = "type assertion failed"
	UnknownReason          = "unknown error found"
	InvalidArgumentReason  = "invalid argument"
	InvalidStateReason     = "invalid state"
	ClientCanceledReason   = "client cancelled"
	TimeoutReason          = "timeout"
	GatewayReason          = "gateway"
	UnauthorizedReason     = "unauthorized" // Not authenticated
	ForbiddenReason        = "forbidden"
	InvalidEventTypeReason = "invalid event type"
	ConditionNotMetReason  = "condition not met"
)

type Error struct {
	Entity     string
	Key        string
	Reason     string
	Details    string
	httpStatus int
}

func (e Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s %s, %s: %s", e.Entity, e.Key, e.Reason, e.Details)
	}
	return fmt.Sprintf("%s %s, %s", e.Entity, e.Reason, e.Key)
}

func (e Error) Unwrap() error {
	return nil
}

func (e Error) Type() string {
	return e.Reason
}

func (e Error) HTTPStatus() int {
	return e.httpStatus
}

func NewUnknownError(entity string, details string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     UnknownReason,
		Details:    details,
		httpStatus: 500,
	}
}

func NewInvalidArgumentError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     InvalidArgumentReason,
		httpStatus: 400,
	}
}

func NewNotFoundError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     NotFoundReason,
		httpStatus: 404,
	}
}

func NewDuplicateError(entity string, details string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     DuplicateReason,
		Details:    details,
		httpStatus: 409,
	}
}

func NewFoundManyError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     FoundManyReason,
		httpStatus: 500,
	}
}

func NewConditionNotMetError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     ConditionNotMetReason,
		httpStatus: 409,
	}
}

func NewTypeAssertionError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     TypeAssertionReason,
		httpStatus: 500,
	}
}

func NewTimeoutError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     TimeoutReason,
		httpStatus: 504,
	}
}

func NewGatewayError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     GatewayReason,
		httpStatus: 502,
	}
}

func NewCancellationError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     ClientCanceledReason,
		httpStatus: 499,
	}
}

func NewInvalidStateError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     InvalidStateReason,
		httpStatus: 409,
	}
}

func NewInvalidEventTypeError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     InvalidEventTypeReason,
		httpStatus: 500,
	}
}

func NewForbiddenError(entity string, keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     entity,
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     ForbiddenReason,
		httpStatus: 403,
	}
}

func NewUnauthorized(keyFmt string, args ...interface{}) Error {
	return Error{
		Entity:     "user",
		Key:        fmt.Sprintf(keyFmt, args...),
		Reason:     UnauthorizedReason,
		httpStatus: 401,
	}
}

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}

	var e Error

	if errors.As(err, &e) {
		return e.Reason == NotFoundReason
	}
	return false
}
