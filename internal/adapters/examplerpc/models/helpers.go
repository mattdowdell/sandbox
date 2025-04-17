package models

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Errors that can be returned from ParseID.
var (
	ErrParseID = errors.New("failed to parse id")
)

// ParseID parses the "Id" field from a Protobuf message into a UUID.
func ParseID(msg interface{ GetId() string }) (uuid.UUID, error) {
	id, err := uuid.Parse(msg.GetId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w %q: %w", ErrParseID, msg.GetId(), err)
	}

	return id, nil
}
