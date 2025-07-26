package datastore

import (
	"embed"
)

// MigrationFS is an embedded filesystem containing all SQL migrations necessary for provisioning
// the database.
//
// Migration filenames must be in the format <version>_<name>.sql, where version is a 4 digit
// number. The first version is 0001, the second 0002, and so on. The leading 0s are stripped from
// the version during migration, so exist only to support numerical ordering when listing files.
//
//go:embed migrations/*.sql
var MigrationFS embed.FS
