package main

import (
	"github.com/go-jet/jet/v2/generator/metadata"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/gofrs/uuid/v5"
)

func updateTemplate() template.Template {
	return template.Default(postgres.Dialect).UseSchema(updateSchema)
}

// updateSchema removes the schema name from the path the generated files are output into.
//
//nolint:gocritic // function signature is defined in a library
func updateSchema(schema metadata.Schema) template.Schema {
	model := template.DefaultModel().UseTable(updateTable)
	return template.DefaultSchema(schema).UseModel(model).UsePath("")
}

// updateTable skips the "goose_db_version" table, used for tracking database migrations.
func updateTable(table metadata.Table) template.TableModel {
	if table.Name == "goose_db_version" {
		return template.TableModel{
			Skip: true,
		}
	}

	return template.DefaultTableModel(table).UseField(updateColumnType)
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
