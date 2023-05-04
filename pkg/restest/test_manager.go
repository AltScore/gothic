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
	useDefaultErrorHandler bool
	errorHandler           echo.HTTPErrorHandler
	handler                echo.HandlerFunc
	body                   []byte
	actualErrorResult      error
	rec                    *httptest.ResponseRecorder
	context                map[string]interface{}
	contextAfterCall       echo.Context
	method                 string
	path                   string
	paramNames             []string
	paramValues            []string
}

type given struct {
	tm *TestManager
}

type when struct {
	tm *TestManager
}

type then struct {
	tm                   *TestManager
	jsonResponseReceived interface{}
}

// For builds a test manager using the provided test context T
func For(t *testing.T) *TestManager {
	return &TestManager{t: t}
}

// WithErrorHandler adds to the TestManager an error handler to convert errors to HTTP codes
func (m *TestManager) WithErrorHandler(errorHandler echo.HTTPErrorHandler) *TestManager {
	m.errorHandler = errorHandler
	return m
}

func (m *TestManager) WithDefaultErrorHandler() *TestManager {
	m.useDefaultErrorHandler = true
	return m
}

// Given returns the GIVEN builder. It allows to define the preconditions for the test
func (m *TestManager) Given() *given {
	return &given{m}
}

// When returns the WHEN builder. It allows to execute the handler under test
func (m *TestManager) When() *when {
	return &when{m}
}

// Then returns the THEN builder. It allows to verify the expectations (post conditions) for the test
func (m *TestManager) Then() *then {
	return &then{tm: m}
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

	for name, value := range m.context {
		c.Set(name, value)
	}

	c.SetParamNames(m.paramNames...)
	c.SetParamValues(m.paramValues...)

	var err error
	if m.handler != nil {
		chain := xapi.ErrorNormalizerMiddleware()(m.handler)

		err = chain(c)
	} else {
		e.ServeHTTP(m.rec, req)
	}
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

// Body sets the body to use when calling the handler
func (g *given) Body(body interface{}) *given {
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
func (g *given) Context(name string, value interface{}) *given {
	if g.tm.context == nil {
		g.tm.context = make(map[string]interface{})
	}
	g.tm.context[name] = value
	return g
}

// PathParam sets a value in the request context under the given name. Used to store entity id, etc.
func (g *given) PathParam(name string, value string) *given {
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
func (g *given) Method(method string) *given {
	g.tm.method = method
	return g
}

// Path sets the path called
func (g *given) Path(path string, args ...interface{}) *given {
	g.tm.path = fmt.Sprintf(path, args...)
	return g
}

// Calls execute the request to the handler. A proper initialized echo.Context is build and used as parameter
func (w *when) Calls(handler func(c echo.Context) error) *when {
	w.tm.handler = handler
	w.tm.exec()
	return w
}

func (w *when) PerformTheCall() *when {
	w.tm.exec()
	return w
}

// StatusCodeIs verifies the response has the expected status code
func (t *then) StatusCodeIs(expectedStatusCode int) *then {
	assert.Equal(t.tm.t, expectedStatusCode, t.tm.rec.Code, "unexpected status code")
	return t
}

// ContentTypeIs verifies the response has the expected content type
func (t *then) ContentTypeIs(expectedType string) *then {
	assert.Contains(t.tm.t, t.tm.rec.Header().Get(HeaderContentType), expectedType)
	return t
}

// ContentTypeIsJson verifies the response has the expected JSON content type
func (t *then) ContentTypeIsJson() *then {
	return t.ContentTypeIs("application/json")
}

// BodyFieldIs verifies the response body contains a filed in the provided path with the expected value.
//
// The path is given in JSONPath format. See https://goessner.net/articles/JsonPath/
func (t *then) BodyFieldIs(jsonPath string, expectedValue interface{}) *then {

	assert.Equal(t.tm.t, expectedValue, t.getJsonFieldValue(jsonPath))

	return t
}

func (t *then) GetContextAfterCall() echo.Context {
	return t.tm.contextAfterCall
}

// BodyIs verifies the response body has the expected value
func (t *then) BodyIs(expectedBody interface{}) *then {
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

func (t *then) getExpectedBytes(expectedBody interface{}) []byte {
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

func (t *then) BodyFieldSatisfies(jsonPath string, predicate func(*testing.T, interface{})) *then {
	predicate(t.tm.t, t.getJsonFieldValue(jsonPath))
	return t
}

func (t *then) BodyFieldHasValue(jsonPath string) *then {
	assert.NotNil(t.tm.t, t.getJsonFieldValue(jsonPath))
	return t
}

// https://pkg.go.dev/github.com/PaesslerAG/jsonpath#example-package--Gval
func (t *then) getJsonFieldValue(jsonPath string) interface{} {
	fieldValue, err := jsonpath.Get(t.fixJsonPath(jsonPath), t.getJson())
	require.Nil(t.tm.t, err)
	return fieldValue
}

func (t *then) fixJsonPath(jsonPath string) string {
	if strings.HasPrefix(jsonPath, "$.") {
		return jsonPath
	}
	return "$." + jsonPath
}

func (t *then) getJson() interface{} {
	if t.jsonResponseReceived == nil {

		t.jsonResponseReceived = interface{}(nil)

		require.Nil(t.tm.t, json.Unmarshal(t.tm.rec.Body.Bytes(), &t.jsonResponseReceived))
	}
	return t.jsonResponseReceived
}

type bytesHolder struct {
	bytes []byte
}
