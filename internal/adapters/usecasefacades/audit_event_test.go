package usecasefacades_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockcommon"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockusecasefacades"
)

func Test_AuditEvent_List(t *testing.T) {
	// arrange
	datastore := mockcommon.NewDatastore(t)

	provider := mockcommon.NewProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := mockusecasefacades.NewAuditEventLister(t)
	usecase.
		EXPECT().
		Execute(t.Context(), datastore).
		Return(nil, nil).
		Once()

	facade := usecasefacades.NewAuditEvent(
		provider,
		usecase,
		mockusecasefacades.NewAuditEventWatcher(t),
	)

	// act
	output, err := facade.List(t.Context())

	// assert
	assert.Empty(t, output)
	assert.NoError(t, err)
}

func Test_AuditEvent_Watch(t *testing.T) {
	// arrange
	datastore := mockcommon.NewDatastore(t)

	provider := mockcommon.NewProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	ch := make(chan *entities.AuditEvent)

	usecase := mockusecasefacades.NewAuditEventWatcher(t)
	usecase.
		EXPECT().
		Execute(t.Context(), datastore).
		Return(ch).
		Once()

	facade := usecasefacades.NewAuditEvent(
		provider,
		mockusecasefacades.NewAuditEventLister(t),
		usecase,
	)

	// act
	output := facade.Watch(t.Context())

	// assert
	assert.NotNil(t, output)
}
