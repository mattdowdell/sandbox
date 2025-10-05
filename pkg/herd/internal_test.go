package herd

// Exports to support unit tests.
const (
	TableNameSystem = tableNameSystem
	TableNameUser   = tableNameUser

	SystemVersionQuery = systemVersionQuery
	UserVersionQuery   = userVersionQuery

	SystemRecordQuery = systemRecordQuery
	UserRecordQuery   = userRecordQuery
)

// Exports to support unit tests.
var (
	NewMigrator = newMigrator
	NewRecorder = newRecorder
)
