package modelhelpers

import (
	gofrsuuid "github.com/gofrs/uuid/v5"
	googleuuid "github.com/google/uuid"
)

func toGofrsUUID(id googleuuid.UUID) gofrsuuid.UUID {
	return gofrsuuid.UUID(id)
}

func toGoogleUUID(id gofrsuuid.UUID) googleuuid.UUID {
	return googleuuid.UUID(id)
}
