package config

import (
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"reflect"
	"strconv"
	"strings"
)

const (
	koanfTag = "koanf"
)

// Sentinel errors that can be returned by Encode.
var ErrNotStruct = errors.New("non-struct found")

var (
	textMarshaler = reflect.TypeFor[*encoding.TextMarshaler]().Elem()
	jsonMarshaler = reflect.TypeFor[*json.Marshaler]().Elem()
)

// Encode converts a configuration struct to a flattened map using the given delimiter for nested
// keys and all values as strings. It is intended for debugging only and is not guaranteed to
// produce an entirely accurate representation of the original configuration.
//
// Encoding all primitive types is supported natively. If a type implements [encoding.TextMarshaler]
// or [encoding/json.Marshaler], that result will be returned instead. Arrays, slices and maps
// cannot currently be encoded, but a placeholder message is returned instead of raising an error.
//
// Private struct fields are ignored.
func Encode[T any](value T, delim string) (map[string]string, error) {
	v := reflect.ValueOf(value)
	return encodeStruct(v, "" /*prefix*/, delim)
}

//nolint:gocognit // no easy way to break up field encoding
func encodeStruct(value reflect.Value, prefix, delim string) (map[string]string, error) {
	if k := value.Kind(); k == reflect.Pointer {
		value = value.Elem()
	}

	if k := value.Kind(); k != reflect.Struct {
		return nil, fmt.Errorf("%w: %s", ErrNotStruct, k)
	}

	result := map[string]string{}

	if prefix != "" {
		encoded, ok, err := marshalFieldValue(value)
		if err != nil {
			return nil, err
		}

		if ok {
			result[prefix] = encoded
			return result, nil
		}
	}

	for i := range value.NumField() {
		val := value.Field(i)
		field := value.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		name := fieldName(&field)

		if val.Kind() == reflect.Struct {
			enc, err := encodeStruct(val, name, delim)
			if err != nil {
				return nil, err
			}

			maps.Copy(result, enc)
			continue
		}

		key := name
		if prefix != "" {
			key = fmt.Sprintf("%s%s%s", prefix, delim, key)
		}

		encoded, err := encodeFieldValue(val)
		if err != nil {
			return nil, err
		}

		result[key] = encoded
	}

	return result, nil
}

func fieldName(field *reflect.StructField) string {
	name := strings.ToLower(field.Name)
	tag := field.Tag.Get(koanfTag)

	// TODO: if anonymous and squash is used, return an empty string

	if n, _, _ := strings.Cut(tag, ","); n != "" {
		name = n
	}

	return name
}

func marshalFieldValue(value reflect.Value) (val string, ok bool, err error) {
	typ := value.Type()

	if typ.Implements(textMarshaler) {
		//nolint:forcetypeassert // already checked via reflection
		m := value.Interface().(encoding.TextMarshaler)

		data, err := m.MarshalText()
		if err != nil {
			return "", true, err
		}

		return string(data), true, nil
	}

	if typ.Implements(jsonMarshaler) {
		//nolint:forcetypeassert // already checked via reflection
		m := value.Interface().(json.Marshaler)

		data, err := m.MarshalJSON()
		if err != nil {
			return "", true, err
		}

		return string(data), true, nil
	}

	return "", false, nil
}

func encodeFieldValue(value reflect.Value) (string, error) {
	encoded, ok, err := marshalFieldValue(value)
	if err != nil {
		return "", err
	}

	if ok {
		return encoded, nil
	}

	//nolint:exhaustive // support just primitive types for now
	switch k := value.Kind(); k {
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10 /*base*/), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10 /*base*/), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'g' /*fmt*/, -1 /*prec*/, 64 /*bitSize*/), nil

	case reflect.String:
		return value.String(), nil

	default:
		return fmt.Sprintf("unsupported kind: %s", k), nil
	}
}
