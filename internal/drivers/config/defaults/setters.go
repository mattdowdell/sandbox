package defaults

import (
	"reflect"
	"strconv"
	"time"
)

func setBool(field reflect.Value, tag string) error {
	val, err := strconv.ParseBool(tag)
	if err != nil {
		return err
	}

	field.SetBool(val)
	return nil
}

func setInt(field reflect.Value, tag string, bitSize int) error {
	val, err := strconv.ParseInt(tag, 0 /*base*/, bitSize)
	if err != nil {
		return err
	}

	field.SetInt(val)
	return nil
}

func setDuration(field reflect.Value, tag string) error {
	val, err := time.ParseDuration(tag)
	if err != nil {
		return err
	}

	field.SetInt(int64(val))
	return nil
}

func setInt64OrDuration(field reflect.Value, tag string) error {
	typ := field.Type()
	if typ.PkgPath() == "time" && typ.Name() == "Duration" {
		return setDuration(field, tag)
	}

	//nolint:mnd // bit size is not particularly magic
	return setInt(field, tag, 64 /*bitSize*/)
}

func setUint(field reflect.Value, tag string, bitSize int) error {
	val, err := strconv.ParseUint(tag, 0 /*base*/, bitSize)
	if err != nil {
		return err
	}

	field.SetUint(val)
	return nil
}

func setFloat(field reflect.Value, tag string, bitSize int) error {
	val, err := strconv.ParseFloat(tag, bitSize)
	if err != nil {
		return err
	}

	field.SetFloat(val)
	return nil
}
