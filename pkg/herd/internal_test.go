package herd

// Exports to support unit tests.
const (
	SystemExistsQuery = systemExistsQuery

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
