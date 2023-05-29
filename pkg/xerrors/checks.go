package xerrors

import (
	"fmt"
	"reflect"
)

// EnsureNotEmpty panics if the given pointer is nil, or value is empty string, or numeric 0.
func EnsureNotEmpty(pointer any, format string, args ...any) {
	if pointer == nil {
		panic(fmt.Sprintf(format, args...))
	}

	switch reflect.TypeOf(pointer).Kind() {
	case reflect.Map, reflect.Array, reflect.Slice:
		if reflect.ValueOf(pointer).IsNil() || reflect.ValueOf(pointer).Len() == 0 {
			panic(fmt.Sprintf(format, args...))
		}
	case reflect.Ptr, reflect.Chan:
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

// EnsureHasKey panics if the given map does not contain the given key.
func EnsureHasKey[K comparable, V any](value map[K]V, key K, format string, args ...any) {
	if _, found := value[key]; !found {
		panic(fmt.Sprintf(format, args...))
	}
}
