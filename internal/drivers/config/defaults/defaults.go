// Package defaults provides support for setting default values in struct fields based on the value
// of the "default" struct tag.
package defaults

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const tagName = "default"

// Set sets the fields of a struct to the values in the "default" struct tag. The target argument
// must be a pointer to a struct.
//
// By default, the following types are natively supported:
//
// - int, int, int16, int32, int64 via [strconv.ParseInt].
// - uint, uint8, uint16, uint32, uint64, uintptr via [strconv.ParseUint].
// - float32, float64 via [strconv.ParseFloat].
// - bool via [strconv.ParseBool].
// - string
// - [time.Duration] via [time.ParseDuration].
//
// Integer default values can be specified in binary, octal, decimal, and hexadecimal by adding the
// appropriate prefix. Unparseable values result in an error being returned.
//
// Structs within structs are also supported under the assumption they contain their own set of
// supported types.
//
// Support for types other than those explicitly supported can implement [encoding.TextUnmarshaler].
// This can also be implemented for type aliases of the supported types to override the default
// behaviour. Types that are not supported natively and do not implement [encoding.TextUmarshaler]
// result in an error being returned.
//
// Defaults are not applied in the following circumstances:
//
// - Private fields are skipped.
// - Fields with a non-zero value are skipped.
// - Empty values of the "default" struct tag are skipped.
func Set(target any) error {
	val := reflect.ValueOf(target)
	if k := val.Kind(); k != reflect.Pointer {
		return fmt.Errorf("expected struct pointer, found: %s", k)
	}

	elem := val.Elem()
	if k := elem.Kind(); k != reflect.Struct {
		return fmt.Errorf("expected struct pointer, found pointer to: %s", k)
	}

	return setStruct(elem)
}

func setStruct(val reflect.Value) error {
	typ := val.Type()
	for i := range typ.NumField() {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		if err := setField(val.Field(i), field.Tag.Get(tagName)); err != nil {
			return fmt.Errorf("set field failed: %s: %w", field.Name, err)
		}
	}

	return nil
}

func setField(field reflect.Value, tag string) error {
	ok, err := shouldSet(field, tag)
	if err != nil {
		return err
	}

	if !ok {
		return nil
	}

	if !isZero(field) {
		return nil
	}

	ok, err = defaultByUnmarshal(field, tag)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return defaultByKind(field, tag)
}

func shouldSet(field reflect.Value, tag string) (bool, error) {
	if !field.CanSet() {
		// if we get here, we probably have a bug
		return false, errors.New("unsettable field")
	}

	if field.Kind() == reflect.Struct {
		return true, nil
	}

	return tag != "", nil
}

func isZero(field reflect.Value) bool {
	return reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface())
}

func defaultByUnmarshal(field reflect.Value, tag string) (bool, error) {
	unmarshaler, ok := field.Addr().Interface().(encoding.TextUnmarshaler)
	if ok && tag != "" {
		if err := unmarshaler.UnmarshalText([]byte(tag)); err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

//nolint:mnd // int sizes are not particularly magic
func defaultByKind(field reflect.Value, tag string) error {
	//nolint:exhaustive // default case is sufficient fallback
	switch k := field.Kind(); k {
	case reflect.Pointer:
		// make sure the pointer is non-nil, otherwise we can't set it
		field.Set(reflect.New(field.Type().Elem()))
		return setField(field.Elem(), tag)

	case reflect.Struct:
		return setStruct(field)

	case reflect.Bool:
		return setBool(field, tag)

	case reflect.Int:
		return setInt(field, tag, strconv.IntSize)

	case reflect.Int8:
		return setInt(field, tag, 8 /*bitSize*/)

	case reflect.Int16:
		return setInt(field, tag, 16 /*bitSize*/)

	case reflect.Int32:
		return setInt(field, tag, 32 /*bitSize*/)

	case reflect.Int64:
		return setInt64OrDuration(field, tag)

	case reflect.Uint, reflect.Uintptr:
		return setUint(field, tag, strconv.IntSize)

	case reflect.Uint8:
		return setUint(field, tag, 8 /*bitSize*/)

	case reflect.Uint16:
		return setUint(field, tag, 16 /*bitSize*/)

	case reflect.Uint32:
		return setUint(field, tag, 32 /*bitSize*/)

	case reflect.Uint64:
		return setUint(field, tag, 64 /*bitSize*/)

	case reflect.Float32:
		return setFloat(field, tag, 32 /*bitSize*/)

	case reflect.Float64:
		return setFloat(field, tag, 64 /*bitSize*/)

	case reflect.String:
		field.SetString(tag)
		return nil

	default:
		return fmt.Errorf("unsupported kind: %s", k)
	}
}
