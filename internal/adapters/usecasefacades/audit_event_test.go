package usecasefacades_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/txn/mocktxn"
	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
)

func Test_AuditEvent_List(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	provider := mocktxn.NewMockProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := usecasefacades.NewMockAuditEventLister(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore).
		Return(nil, nil).
		Once()

	facade := usecasefacades.NewAuditEvent(
		provider,
		usecase,
		usecasefacades.NewMockAuditEventWatcher(t),
	)

	// act
	output, err := facade.List(t.Context(), logger)

	// assert
	assert.Empty(t, output)
	assert.NoError(t, err)
}

func Test_AuditEvent_Watch(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	provider := mocktxn.NewMockProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	ch := make(chan *entities.AuditEvent)

	usecase := usecasefacades.NewMockAuditEventWatcher(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore).
		Return(ch).
		Once()

	facade := usecasefacades.NewAuditEvent(
		provider,
		usecasefacades.NewMockAuditEventLister(t),
		usecase,
	)

	// act
	output := facade.Watch(t.Context(), logger)

	// assert
	assert.NotNil(t, output)
}
