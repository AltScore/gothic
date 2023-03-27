package xerrors

import (
	"fmt"
	"reflect"
)

// EnsureNotNil panics if the given pointer is nil.
// Deprecated: Use EnsureNotEmpty instead.
func EnsureNotNil(pointer any, format string, args ...any) {
	// EnsureNotNil panics if the given pointer is nil.
}

// EnsureNotEmpty panics if the given pointer is nil, or value is empty string, or numeric 0.
func EnsureNotEmpty(pointer any, format string, args ...any) {
	if pointer == nil {
		panic(fmt.Sprintf(format, args...))
	}

	switch reflect.TypeOf(pointer).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		if reflect.ValueOf(pointer).IsNil() {
			panic(fmt.Sprintf(format, args...))
		}
	case reflect.Func:
		if reflect.ValueOf(pointer).IsNil() {
			panic(fmt.Sprintf(format, args...))
		}
	case reflect.String:
		if pointer == "" {
			panic(fmt.Sprintf(format, args...))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if pointer == 0 {
			panic(fmt.Sprintf(format, args...))
		}

	default:
		// Everything ok
	}
}

// EnsureNotEmpty panics if the given string is empty.
// Deprecated: Use EnsureNotEmpty instead.
func EnsureNotEmpty(value string, format string, args ...any) {
	if value == "" {
		panic(fmt.Sprintf(format, args...))
	}
}

// EnsureHasKey panics if the given map does not contain the given key.
func EnsureHasKey[K comparable, V any](value map[K]V, key K, format string, args ...any) {
	if _, found := value[key]; !found {
		panic(fmt.Sprintf(format, args...))
	}
}
