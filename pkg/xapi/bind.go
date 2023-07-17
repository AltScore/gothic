package xapi

import (
	"github.com/labstack/echo/v4"
)

// Validating is an interface that can be used to validate a struct.
type Validating interface {
	Validate() error
}

// BindValidated binds the request (body/param/query) to the given struct and validates it.
// If the struct is invalid, it returns an error.
func BindValidated[T Validating](c echo.Context) (T, error) {
	var t T

	if err := c.Bind(&t); err != nil {
		return t, err
	}

	if err := t.Validate(); err != nil {
		return t, err
	}

	return t, nil
}
