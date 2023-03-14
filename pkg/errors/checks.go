package errors

import (
	"fmt"
	"reflect"
)

// EnsureNotNil panics if the given pointer is nil.
func EnsureNotNil(pointer any, format string, args ...any) {
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
	default:
		// Everything ok
	}
}

// EnsureNotEmpty panics if the given string is empty.
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
