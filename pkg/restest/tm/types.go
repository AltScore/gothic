package tm

import (
	"testing"

	"github.com/AltScore/gothic/v2/pkg/restest"
)

// O is a type to create a literal JSON object
//
//	jsonLiteral := tm.O{"name": "a name", "age": 32}
type O = map[string]interface{}

// A is a type to create a literal JSON array
//
//	jsonLiteral := tm.A{ "a vale", "other value" }
type A = []interface{}

// For builds a test manager using the provided test context T
func For(t *testing.T) *restest.TestManager {
	return restest.For(t)
}
