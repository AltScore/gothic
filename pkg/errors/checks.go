package errors

import (
	"fmt"
	"reflect"
)

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
