package xvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AltScore/money/pkg/money"
	"github.com/AltScore/money/pkg/percent"
	"github.com/go-playground/validator/v10"
)

var (
	validate    *validator.Validate
	timeType    = reflect.TypeOf(time.Time{})
	moneyType   = reflect.TypeOf(money.Money{})
	percentType = reflect.TypeOf(percent.Percent(0))
)

func init() {
	validate = validator.New()

	_ = validate.RegisterValidation("noEmpty", noEmpty)
	_ = validate.RegisterValidation("vgte", isGte)
	_ = validate.RegisterValidation("vgt", isGt)
	_ = validate.RegisterValidation("minLen", minLen)
	_ = validate.RegisterValidation("json_array", isJSONArray)
}

func Instance() *validator.Validate {
	return validate
}

// Struct validates a structs exposed fields, and automatically validates nested structs, unless otherwise specified.
//
// It returns InvalidValidationError for bad values passed in and nil or ValidationErrors as error otherwise.
// You will need to assert the error if it's not nil eg. err.(validator.ValidationErrors) to access the array of errors.
func Struct(s interface{}) error {
	return validate.Struct(s)
}

func noEmpty(fl validator.FieldLevel) bool {
	field := fl.Field().String()

	return len(strings.TrimSpace(field)) > 0
}

// isGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asFloat(param)

		return asFloat(field.String()) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt64(param)

		return int64(field.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt64(param)

		return field.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() >= p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {

			now := time.Now().UTC()
			t := field.Convert(timeType).Interface().(time.Time) //nolint:forcetypeassert // Already checked for timeType

			return t.After(now) || t.Equal(now)
		} else if field.Type().ConvertibleTo(moneyType) {
			p := asFloat(param)

			return field.Convert(moneyType).Interface().(money.Money).Number() >= p
		} else if field.Type().ConvertibleTo(percentType) {
			p := asFloat(param)

			return field.Convert(percentType).Interface().(percent.Percent).Number() >= p
		}

	default: // Nothing to do
	}
	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// isGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGt(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asFloat(param)

		return asFloat(field.String()) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt64(param)

		return int64(field.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt64(param)

		return field.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)

		return field.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)

		return field.Float() > p

	case reflect.Struct:

		if field.Type().ConvertibleTo(timeType) {

			now := time.Now().UTC()
			t := field.Convert(timeType).Interface().(time.Time) //nolint:forcetypeassert // Already checked for timeType

			return t.After(now)
		} else if field.Type().ConvertibleTo(moneyType) {
			p := asFloat(param)

			return field.Convert(moneyType).Interface().(money.Money).Number() > p
		} else if field.Type().ConvertibleTo(percentType) {
			p := asFloat(param)

			return field.Convert(percentType).Interface().(percent.Percent).Number() > p
		}
	default: // Nothing to do
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

// asInt64 returns the parameter as an int64
// or panics if it can't convert
func asInt64(param string) int64 {
	return orPanic(strconv.ParseInt(param, 0, 64))
}

// asUint returns the parameter as an uint64
// or panics if it can't convert
func asUint(param string) uint64 {
	return orPanic(strconv.ParseUint(param, 0, 64))
}

// asFloat returns the parameter as a float64
// or panics if it can't convert
func asFloat(param string) float64 {
	return orPanic(strconv.ParseFloat(param, 64))
}

func orPanic[T any](t T, err error) T {
	if err != nil {
		panic(err.Error())
	}

	return t
}

// minLen is the validation function for validating if the current field's length is greater than or equal to the param's value.
func minLen(fl validator.FieldLevel) bool {
	field := fl.Field()
	param := fl.Param()

	switch field.Kind() {

	case reflect.String:
		p := asInt64(param)

		return int64(len(field.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt64(param)

		return int64(field.Len()) >= p

	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
}

var arrayStartRegex = regexp.MustCompile(`^\s*\[`)

// isJSONArray is the validation function for validating if the current field's value is a valid json array string.
func isJSONArray(fl validator.FieldLevel) bool {
	field := fl.Field()

	if field.Kind() == reflect.String {
		val := field.String()
		if valid := json.Valid([]byte(val)); !valid {
			return false
		}

		// If the string is a valid json array, it must start with a '['
		return arrayStartRegex.MatchString(val)
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
