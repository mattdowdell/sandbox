package herdconv

import "go.opentelemetry.io/otel/attribute"

const (
	HerdTableNameKey     = attribute.Key("herd.table.name")
	HerdVersionBeforeKey = attribute.Key("herd.version.before")
	HerdVersionAfterKey  = attribute.Key("herd.version.after")
)

func HerdTableName(value string) attribute.KeyValue {
	return HerdTableNameKey.String(value)
}

func HerdVersionBefore(value int) attribute.KeyValue {
	return HerdVersionBeforeKey.Int(value)
}

func HerdVersionAfter(value int) attribute.KeyValue {
	return HerdVersionAfterKey.Int(value)
}
