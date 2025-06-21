package modelhelpers

import (
	"testing"

	gofrsuuid "github.com/gofrs/uuid/v5"
	googleuuid "github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func ToGoogleUUID(id gofrsuuid.UUID) googleuuid.UUID {
	return toGoogleUUID(id)
}

func Test_toGoogleUUID(t *testing.T) {
	testCases := []struct {
		name string
		have gofrsuuid.UUID
	}{
		{
			name: "v4",
			have: gofrsuuid.Must(gofrsuuid.NewV4()),
		},
		{
			name: "v7",
			have: gofrsuuid.Must(gofrsuuid.NewV7()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			got := toGoogleUUID(tc.have)

			// assert
			assert.Equal(t, tc.have.String(), got.String())
		})
	}
}

func Test_toGofrsUUID(t *testing.T) {
	testCases := []struct {
		name string
		have googleuuid.UUID
	}{
		{
			name: "v4",
			have: googleuuid.Must(googleuuid.NewRandom()),
		},
		{
			name: "v7",
			have: googleuuid.Must(googleuuid.NewV7()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// arrange

			// act
			got := toGofrsUUID(tc.have)

			// assert
			assert.Equal(t, tc.have.String(), got.String())
		})
	}
}
