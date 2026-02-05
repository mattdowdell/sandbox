package herdconv

import "go.opentelemetry.io/otel/attribute"

const (
	HerdMigrationTypeKey = attribute.Key("herd.migration.type")
	HerdVersionBeforeKey = attribute.Key("herd.version.before")
	HerdVersionAfterKey  = attribute.Key("herd.version.after")
)

var (
	HerdMigrationTypeSystem = HerdMigrationTypeKey.String("system")
	HerdMigrationTypeUser   = HerdMigrationTypeKey.String("user")
)

func HerdVersionBefore(value int) attribute.KeyValue {
	return HerdVersionBeforeKey.Int(value)
}

func HerdVersionAfter(value int) attribute.KeyValue {
	return HerdVersionAfterKey.Int(value)
}
