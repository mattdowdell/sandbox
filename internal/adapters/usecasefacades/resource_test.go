package usecasefacades_test

import (
	"log/slog"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/txn/mocktxn"
	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/internal/domain/repositories"
)

const (
	testLimit = 50
)

func Test_Resource_Create(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	ender := mocktxn.NewMockEnder(t)
	ender.EXPECT().Commit().Return(nil).Once()
	ender.EXPECT().Rollback().Return(nil).Once()

	provider := mocktxn.NewMockProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, ender, nil).
		Once()

	usecase := usecasefacades.NewMockResourceCreator(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, &entities.Resource{}).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecase,
		usecasefacades.NewMockResourceGetter(t),
		usecasefacades.NewMockResourceLister(t),
		usecasefacades.NewMockResourceUpdater(t),
		usecasefacades.NewMockResourceDeleter(t),
	)

	// act
	output, err := facade.Create(t.Context(), logger, &entities.Resource{})

	// assert
	assert.NotNil(t, output)
	assert.NoError(t, err)
}

func Test_Resource_Get(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())

	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	provider := mocktxn.NewMockProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := usecasefacades.NewMockResourceGetter(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, id).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecasefacades.NewMockResourceCreator(t),
		usecase,
		usecasefacades.NewMockResourceLister(t),
		usecasefacades.NewMockResourceUpdater(t),
		usecasefacades.NewMockResourceDeleter(t),
	)

	// act
	output, err := facade.Get(t.Context(), logger, id)

	// assert
	assert.NotNil(t, output)
	assert.NoError(t, err)
}

func Test_Resource_List(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	pager := repositories.Pager{
		Limit: testLimit,
	}

	provider := mocktxn.NewMockProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := usecasefacades.NewMockResourceLister(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, pager).
		Return(&repositories.Paged[*entities.Resource]{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecasefacades.NewMockResourceCreator(t),
		usecasefacades.NewMockResourceGetter(t),
		usecase,
		usecasefacades.NewMockResourceUpdater(t),
		usecasefacades.NewMockResourceDeleter(t),
	)

	// act
	output, err := facade.List(t.Context(), logger, pager)

	// assert
	assert.Empty(t, output)
	assert.NoError(t, err)
}

func Test_Resource_Update(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)
	ender := mocktxn.NewMockEnder(t)
	ender.EXPECT().Commit().Return(nil).Once()
	ender.EXPECT().Rollback().Return(nil).Once()

	provider := mocktxn.NewMockProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, ender, nil).
		Once()

	usecase := usecasefacades.NewMockResourceUpdater(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, &entities.Resource{}).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecasefacades.NewMockResourceCreator(t),
		usecasefacades.NewMockResourceGetter(t),
		usecasefacades.NewMockResourceLister(t),
		usecase,
		usecasefacades.NewMockResourceDeleter(t),
	)

	// act
	output, err := facade.Update(t.Context(), logger, &entities.Resource{})

	// assert
	assert.NotNil(t, output)
	assert.NoError(t, err)
}

func Test_Resource_Delete(t *testing.T) {
	// arrange
	id := uuid.Must(uuid.NewV7())

	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewMockDatastore(t)

	ender := mocktxn.NewMockEnder(t)
	ender.EXPECT().Commit().Return(nil).Once()
	ender.EXPECT().Rollback().Return(nil).Once()

	provider := mocktxn.NewMockProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, ender, nil).
		Once()

	usecase := usecasefacades.NewMockResourceDeleter(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, id).
		Return(nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecasefacades.NewMockResourceCreator(t),
		usecasefacades.NewMockResourceGetter(t),
		usecasefacades.NewMockResourceLister(t),
		usecasefacades.NewMockResourceUpdater(t),
		usecase,
	)

	// act
	err := facade.Delete(t.Context(), logger, id)

	// assert
	assert.NoError(t, err)
}
