package xapi

import (
	"bytes"
	"github.com/AltScore/gothic/v2/pkg/xvalidator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

type SomeTypePtr struct {
	Required string `param:"id" validate:"required"`
}

func (s *SomeTypePtr) Validate() error {
	return xvalidator.Struct(s)
}

type SomeType struct {
	Required string `param:"id" validate:"required"`
}

func (s SomeType) Validate() error {
	return xvalidator.Struct(s)
}

func TestBindValidates_works_with_pointers(t *testing.T) {

	c := newContext("POST", "/segment1/12345/segment3", "{}")
	tt, err := BindValidated[*SomeTypePtr](c)

	require.NoError(t, err)
	require.Equal(t, "12345", tt.Required)
}

func TestBindValidates_works_with_structs(t *testing.T) {

	c := newContext("POST", "/segment1/12345/segment3", "{}")
	tt, err := BindValidated[SomeType](c)

	require.NoError(t, err)
	require.Equal(t, "12345", tt.Required)
}

func newContext(method string, path string, body string) echo.Context {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(body)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()

	e := echo.New()
	e.Add("POST", "/segment1/:id/segment3", func(c echo.Context) error {
		return nil
	})

	c := e.NewContext(req, res)

	e.Router().Find(method, path, c)

	return c
}
