package main

import (
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/gofrs/uuid/v5"

	"github.com/mattdowdell/sandbox/pkg/herd"
)

func updateTemplate() template.Template {
	return template.Default(postgres.Dialect).UseSchema(updateSchema)
}

// updateSchema removes the schema name from the path the generated files are output into.
//
//nolint:gocritic // function signature is defined in a library
func updateSchema(schema metadata.Schema) template.Schema {
	model := template.DefaultModel().UseTable(updateModelTable)
	builder := template.DefaultSQLBuilder().UseTable(updateBuilderTable)

	return template.DefaultSchema(schema).UseModel(model).UseSQLBuilder(builder).UsePath("")
}

// updateModelTable skips the "herd_system_migrations" and "herd_user_migrations" tables, which are
// used for tracking database migrations.
func updateModelTable(table metadata.Table) template.TableModel {
	switch table.Name {
	case herd.TableNameSystem, herd.TableNameUser:
		return template.TableModel{
			Skip: true,
		}

	default:
		return template.DefaultTableModel(table).UseField(updateColumnType)
	}
}

// updateBuilderTable skips the "herd_system_migrations" and "herd_user_migrations" tables, which
// are used for tracking database migrations.
func updateBuilderTable(table metadata.Table) template.TableSQLBuilder {
	switch table.Name {
	case herd.TableNameSystem, herd.TableNameUser:
		return template.TableSQLBuilder{
			Skip: true,
		}

	default:
		return template.DefaultTableSQLBuilder(table)
	}
}

// updateColumnType changes the type of UUID columns from github.com/google/uuid to
// github.com/gofrs/uuid/v5.
//
//nolint:gocritic // function signature is defined in a library
func updateColumnType(column metadata.Column) template.TableModelField {
	field := template.DefaultTableModelField(column)

	if column.DataType.Name == "uuid" {
		field.Type = template.NewType(uuid.UUID{})
	}

	return field
}
