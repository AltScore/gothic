package restest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/AltScore/gothic/pkg/xapi"
	"github.com/PaesslerAG/jsonpath"
	"github.com/labstack/echo/v4"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/totemcaf/gollections/slices"
)

const HeaderContentType = "Content-Type"

// TestManager is helper to test REST handler created with echo HTTP server
type TestManager struct {
	t                      *testing.T
	server                 *echo.Echo
	module                 xapi.Module // this is the subject under test
	useDefaultErrorHandler bool
	errorHandler           echo.HTTPErrorHandler
	handler                echo.HandlerFunc // Deprecated
	body                   []byte
	actualErrorResult      error
	rec                    *httptest.ResponseRecorder
	context                map[string]interface{}
	contextAfterCall       echo.Context
	method                 string
	path                   string
	paramNames             []string // Deprected
	paramValues            []string // Deprecated
}

type GivenWrapper struct {
	tm *TestManager
}

type WhenWrapper struct {
	tm *TestManager
}

type ThenWrapper struct {
	tm                   *TestManager
	jsonResponseReceived interface{}
}

// For builds a test manager using the provided test context T
func For(t *testing.T) *TestManager {
	return &TestManager{t: t, useDefaultErrorHandler: true}
}

// WithErrorHandler adds to the TestManager an error handler to convert errors to HTTP codes
func (m *TestManager) WithErrorHandler(errorHandler echo.HTTPErrorHandler) *TestManager {
	m.errorHandler = errorHandler
	return m
}

// WithDefaultErrorHandler adds to the TestManager the default error handler to convert errors to HTTP codes
// This is the default behavior. TO remove the default error handler, use WithErrorHandler(false)
func (m *TestManager) WithDefaultErrorHandler(useDefault ...bool) *TestManager {
	m.useDefaultErrorHandler = len(useDefault) == 0 || useDefault[0]
	return m
}

// Given returns the GIVEN builder. It allows to define the preconditions for the test
func (m *TestManager) Given() *GivenWrapper {
	return &GivenWrapper{m}
}

// When returns the WHEN builder. It allows to execute the handler under test
func (m *TestManager) When() *WhenWrapper {
	return &WhenWrapper{m}
}

// Then returns the THEN builder. It allows to verify the expectations (post conditions) for the test
func (m *TestManager) Then() *ThenWrapper {
	return &ThenWrapper{tm: m}
}

// Json allows to define the request or response body from JSON in a string.
func (m *TestManager) Json(jsonStr string) interface{} {
	return bytesHolder{[]byte(jsonStr)}
}

// JsonFile allows to define the request or response body from JSON in a string.
func (m *TestManager) JsonFile(fileName string) interface{} {
	jsonBytes, err := os.ReadFile(fileName)

	m.addErr(err)

	return bytesHolder{jsonBytes}
}

func (m *TestManager) addErr(err error) {
	if err != nil {
		m.t.Error(err)
	}
}

// See https://github.com/labstack/echo/discussions/2098
func (m *TestManager) exec() {
	e := m.server

	if e == nil {
		e = echo.New()
	}

	req := httptest.NewRequest(m.method, m.path, bytes.NewReader(m.body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	m.rec = httptest.NewRecorder()

	c := e.NewContext(req, m.rec)

	m.setRoutes(e)

	for name, value := range m.context {
		c.Set(name, value)
	}

	// Deprecated, use Module
	c.SetParamNames(m.paramNames...)
	c.SetParamValues(m.paramValues...)

	// find the handler for the given path

	if m.handler == nil {
		e.Router().Find(m.method, m.path, c)
		m.handler = c.Handler()
	}

	// Apply middlewares
	chain := xapi.ErrorNormalizerMiddleware()(m.handler)

	err := chain(c)

	if err != nil {
		if m.errorHandler != nil {
			m.errorHandler(err, c)
		} else if m.useDefaultErrorHandler {
			e.DefaultHTTPErrorHandler(err, c)
		}
	}

	m.contextAfterCall = c
	m.actualErrorResult = err
}

func (m *TestManager) WithServer(server *echo.Echo) *TestManager {
	m.server = server
	return m
}

// setRoutes loads the existent routes in module into the router
func (m *TestManager) setRoutes(e *echo.Echo) {
	if m.module == nil {
		return // backward compatibility
	}

	for _, route := range m.module.Routes() {
		e.Add(route.Method, route.Path, route.Handler)
	}
}

// Body sets the body to use when calling the handler
func (g *GivenWrapper) Body(body interface{}) *GivenWrapper {
	if bh, ok := body.(bytesHolder); ok {
		g.tm.body = bh.bytes
	} else {
		bodyBytes, err := json.Marshal(body)
		g.tm.addErr(err)
		g.tm.body = bodyBytes
	}

	return g
}

// Context sets a value in the request context under the given name. Used to store current user, current request id, etc.
func (g *GivenWrapper) Context(name string, value interface{}) *GivenWrapper {
	if g.tm.context == nil {
		g.tm.context = make(map[string]interface{})
	}
	g.tm.context[name] = value
	return g
}

// PathParam sets a value in the request context under the given name. Used to store entity id, etc.
// Deprecated, use Module() and CallPath() instead
func (g *GivenWrapper) PathParam(name string, value string) *GivenWrapper {

	idx, found := slices.Index2(g.tm.paramNames, name)
	if found {
		g.tm.paramValues[idx] = value
	} else {
		g.tm.paramNames = append(g.tm.paramNames, name)
		g.tm.paramValues = append(g.tm.paramValues, value)
	}

	return g
}

// Method sets the method to use
func (g *GivenWrapper) Method(method string) *GivenWrapper {
	g.tm.method = method
	return g
}

// Path sets the path called
func (g *GivenWrapper) Path(path string, args ...interface{}) *GivenWrapper {
	g.tm.path = fmt.Sprintf(path, args...)
	return g
}

func (g *GivenWrapper) Module(module xapi.Module) *GivenWrapper {
	g.tm.module = module
	return g
}

// Calls execute the request to the handler. A proper initialized echo.Context is build and used as parameter
// Deprecated, use Module() and CallPath() instead
func (w *WhenWrapper) Calls(handler func(c echo.Context) error) *WhenWrapper {
	w.tm.handler = handler
	w.tm.exec()
	return w
}

// PerformTheCall
// Deprecated, use Module() and CallPath() instead
func (w *WhenWrapper) PerformTheCall() *WhenWrapper {
	w.tm.exec()
	return w
}

func (w *WhenWrapper) CallsPath(path string) *WhenWrapper {
	w.tm.path = path
	w.tm.exec()
	return w
}

// StatusCodeIs verifies the response has the expected status code
func (t *ThenWrapper) StatusCodeIs(expectedStatusCode int) *ThenWrapper {
	assert.Equal(t.tm.t, expectedStatusCode, t.tm.rec.Code, "unexpected status code")
	return t
}

// ContentTypeIs verifies the response has the expected content type
func (t *ThenWrapper) ContentTypeIs(expectedType string) *ThenWrapper {
	assert.Contains(t.tm.t, t.tm.rec.Header().Get(HeaderContentType), expectedType)
	return t
}

// ContentTypeIsJson verifies the response has the expected JSON content type
func (t *ThenWrapper) ContentTypeIsJson() *ThenWrapper {
	return t.ContentTypeIs("application/json")
}

// BodyFieldIs verifies the response body contains a filed in the provided path with the expected value.
//
// The path is given in JSONPath format. See https://goessner.net/articles/JsonPath/
func (t *ThenWrapper) BodyFieldIs(jsonPath string, expectedValue interface{}) *ThenWrapper {

	assert.Equal(t.tm.t, expectedValue, t.getJsonFieldValue(jsonPath))

	return t
}

func (t *ThenWrapper) GetContextAfterCall() echo.Context {
	return t.tm.contextAfterCall
}

// BodyIs verifies the response body has the expected value
func (t *ThenWrapper) BodyIs(expectedBody interface{}) *ThenWrapper {
	expectedBodyBytes := t.getExpectedBytes(expectedBody)

	actualBytes := t.tm.rec.Body.Bytes()

	options := jsondiff.DefaultConsoleOptions()
	options.SkipMatches = true

	differences, explanation := jsondiff.Compare(expectedBodyBytes, actualBytes, &options)

	if differences != jsondiff.FullMatch {
		assert.Fail(t.tm.t, explanation)
	}

	return t
}

func (t *ThenWrapper) getExpectedBytes(expectedBody interface{}) []byte {
	var expectedBodyBytes []byte
	if bh, ok := expectedBody.(bytesHolder); ok {
		expectedBodyBytes = bh.bytes
	} else {
		var err error
		expectedBodyBytes, err = json.Marshal(expectedBody)

		t.tm.addErr(err)
	}
	return expectedBodyBytes
}

func (t *ThenWrapper) BodyFieldSatisfies(jsonPath string, predicate func(*testing.T, interface{})) *ThenWrapper {
	predicate(t.tm.t, t.getJsonFieldValue(jsonPath))
	return t
}

func (t *ThenWrapper) BodyFieldHasValue(jsonPath string) *ThenWrapper {
	assert.NotNil(t.tm.t, t.getJsonFieldValue(jsonPath))
	return t
}

// https://pkg.go.dev/github.com/PaesslerAG/jsonpath#example-package--Gval
func (t *ThenWrapper) getJsonFieldValue(jsonPath string) interface{} {
	fieldValue, err := jsonpath.Get(t.fixJsonPath(jsonPath), t.getJson())
	require.Nil(t.tm.t, err)
	return fieldValue
}

func (t *ThenWrapper) fixJsonPath(jsonPath string) string {
	if strings.HasPrefix(jsonPath, "$.") {
		return jsonPath
	}
	return "$." + jsonPath
}

func (t *ThenWrapper) getJson() interface{} {
	if t.jsonResponseReceived == nil {

		t.jsonResponseReceived = interface{}(nil)

		require.Nil(t.tm.t, json.Unmarshal(t.tm.rec.Body.Bytes(), &t.jsonResponseReceived))
	}
	return t.jsonResponseReceived
}

type bytesHolder struct {
	bytes []byte
}
