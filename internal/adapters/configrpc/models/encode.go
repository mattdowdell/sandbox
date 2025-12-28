package models

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

var ErrNotStruct = errors.New("non-struct found")

var (
	textMarshaler = reflect.TypeOf(new(encoding.TextMarshaler)).Elem()
	jsonMarshaler = reflect.TypeOf(new(json.Marshaler)).Elem()
)

// Encode converts a configuration struct to a flattened map.
func Encode[T any](value T, delim string) (map[string]string, error) {
	v := reflect.ValueOf(value)
	return encodeStruct(v, "" /*prefix*/, delim)
}

// ...
func encodeStruct(value reflect.Value, prefix, delim string) (map[string]string, error) {
	if k := value.Kind(); k == reflect.Pointer {
		value = value.Elem()
	}

	if k := value.Kind(); k != reflect.Struct {
		return nil, fmt.Errorf("%w: %s", ErrNotStruct, k)
	}

	result := map[string]string{}

	// TODO: check for marshal implementation, e.g. time.Time

	for i := range value.NumField() {
		val := value.Field(i)
		field := value.Type().Field(i)
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

		if isRedacted(&field) {
			result[key] = "********"
			continue
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
	tag := field.Tag.Get("koanf")

	// TODO: if anonymous and squash is used, return an empty string

	if n, _, _ := strings.Cut(tag, ","); n != "" {
		name = n
	}

	return name
}

func isRedacted(field *reflect.StructField) bool {
	tag := field.Tag.Get("json")

	n, _, _ := strings.Cut(tag, ",")
	return n == "-"
}

//nolint:revive // boolean last makes more sense in this case
func marshalFieldValue(value reflect.Value) (string, error, bool) {
	typ := value.Type()

	if typ.Implements(textMarshaler) {
		//nolint:forcetypeassert // already checked via reflection
		m := value.Interface().(encoding.TextMarshaler)

		data, err := m.MarshalText()
		if err != nil {
			return "", err, true
		}

		return string(data), nil, true
	}

	if typ.Implements(jsonMarshaler) {
		//nolint:forcetypeassert // already checked via reflection
		m := value.Interface().(json.Marshaler)

		data, err := m.MarshalJSON()
		if err != nil {
			return "", err, true
		}

		return string(data), nil, true
	}

	return "", nil, false
}

func encodeFieldValue(value reflect.Value) (string, error) {
	encoded, err, ok := marshalFieldValue(value)
	if ok {
		return encoded, err
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
