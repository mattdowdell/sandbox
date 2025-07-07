package usecasefacades_test

import (
	"log/slog"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/stretchr/testify/assert"

	"github.com/mattdowdell/sandbox/internal/adapters/usecasefacades"
	"github.com/mattdowdell/sandbox/internal/domain/entities"
	"github.com/mattdowdell/sandbox/mocks/adapters/mocktxn"
	"github.com/mattdowdell/sandbox/mocks/adapters/mockusecasefacades"
)

func Test_Resource_Create(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewDatastore(t)

	commit := mocktxn.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mocktxn.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mocktxn.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	usecase := mockusecasefacades.NewResourceCreator(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, &entities.Resource{}).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		usecase,
		mockusecasefacades.NewResourceGetter(t),
		mockusecasefacades.NewResourceLister(t),
		mockusecasefacades.NewResourceUpdater(t),
		mockusecasefacades.NewResourceDeleter(t),
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
	datastore := mocktxn.NewDatastore(t)

	provider := mocktxn.NewProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := mockusecasefacades.NewResourceGetter(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, id).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		mockusecasefacades.NewResourceCreator(t),
		usecase,
		mockusecasefacades.NewResourceLister(t),
		mockusecasefacades.NewResourceUpdater(t),
		mockusecasefacades.NewResourceDeleter(t),
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
	datastore := mocktxn.NewDatastore(t)

	provider := mocktxn.NewProvider(t)
	provider.EXPECT().Datastore().Return(datastore).Once()

	usecase := mockusecasefacades.NewResourceLister(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore).
		Return(nil, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		mockusecasefacades.NewResourceCreator(t),
		mockusecasefacades.NewResourceGetter(t),
		usecase,
		mockusecasefacades.NewResourceUpdater(t),
		mockusecasefacades.NewResourceDeleter(t),
	)

	// act
	output, err := facade.List(t.Context(), logger)

	// assert
	assert.Empty(t, output)
	assert.NoError(t, err)
}

func Test_Resource_Update(t *testing.T) {
	// arrange
	logger := slog.New(slog.DiscardHandler)
	datastore := mocktxn.NewDatastore(t)

	commit := mocktxn.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mocktxn.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mocktxn.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	usecase := mockusecasefacades.NewResourceUpdater(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, &entities.Resource{}).
		Return(&entities.Resource{}, nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		mockusecasefacades.NewResourceCreator(t),
		mockusecasefacades.NewResourceGetter(t),
		mockusecasefacades.NewResourceLister(t),
		usecase,
		mockusecasefacades.NewResourceDeleter(t),
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
	datastore := mocktxn.NewDatastore(t)

	commit := mocktxn.NewCommitFn(t)
	commit.EXPECT().Execute().Return(nil).Once()

	rollback := mocktxn.NewRollbackFn(t)
	rollback.EXPECT().Execute().Return(nil).Once()

	provider := mocktxn.NewProvider(t)
	provider.
		EXPECT().
		BeginTx(t.Context()).
		Return(datastore, commit.Execute, rollback.Execute, nil).
		Once()

	usecase := mockusecasefacades.NewResourceDeleter(t)
	usecase.
		EXPECT().
		Execute(t.Context(), logger, datastore, id).
		Return(nil).
		Once()

	facade := usecasefacades.NewResource(
		provider,
		mockusecasefacades.NewResourceCreator(t),
		mockusecasefacades.NewResourceGetter(t),
		mockusecasefacades.NewResourceLister(t),
		mockusecasefacades.NewResourceUpdater(t),
		usecase,
	)

	// act
	err := facade.Delete(t.Context(), logger, id)

	// assert
	assert.NoError(t, err)
}
