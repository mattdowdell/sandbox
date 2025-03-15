package workerpool_test

import (
	"context"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/mattdowdell/sandbox/internal/drivers/workerpool"
	"github.com/mattdowdell/sandbox/mocks/drivers/mockworkerpool"
)

const (
	testSize = 2
)

func Test_New_Success(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	// act
	pool, err := workerpool.New(testSize, handler, collector)

	// assert
	assert.NotNil(t, pool)
	assert.NoError(t, err)
}

func Test_New_Error(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	// act
	pool, err := workerpool.New(0 /*size*/, handler, collector)

	// assert
	assert.Nil(t, pool)
	assert.EqualError(t, err, "size must be > 0")
}

func Test_Pool_Start(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	pool, err := workerpool.New(testSize, handler, collector)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())

	// act
	go func() {
		pool.Start(ctx)
	}()

	<-pool.Started()
	cancel()

	// assert
	<-pool.Complete()
}

func Test_Pool_Wait_Success(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	pool, err := workerpool.New(testSize, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()

	// act
	assert.NoError(t, pool.Wait(t.Context()))

	// no assert necessary
}

func Test_Pool_Wait_Error(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	pool, err := workerpool.New(testSize, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	// act
	err = pool.Wait(ctx)

	// assert
	assert.EqualError(t, err, "failed to wait for workers: context canceled")
}

func Test_Pool_Add_Success(t *testing.T) {
	// arrange
	var results []int

	handler := mockworkerpool.NewHandler[int, int](t)
	handler.
		EXPECT().
		Handle(mock.AnythingOfType("*context.cancelCtx"), mock.AnythingOfType("int")).
		RunAndReturn(func(_ context.Context, value int) int {
			return value * 2
		}).
		Times(5)

	collector := mockworkerpool.NewCollector[int](t)
	collector.
		EXPECT().
		Collect(mock.AnythingOfType("int")).
		RunAndReturn(func(value int) {
			results = append(results, value)
		}).
		Times(5)

	pool, err := workerpool.New(testSize, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()

	// act
	for i := range 5 {
		require.NoError(t, pool.Add(t.Context(), i))
	}

	// assert
	require.NoError(t, pool.Wait(t.Context()))

	slices.Sort(results)
	assert.Equal(t, []int{0, 2, 4, 6, 8}, results)
}

func Test_Pool_Add_Error(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	collector := mockworkerpool.NewCollector[int](t)

	pool, err := workerpool.New(testSize, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()
	require.NoError(t, pool.Wait(t.Context()))

	// act
	err = pool.Add(t.Context(), 1)

	// assert
	assert.EqualError(t, err, "worker pool queue has been closed")
}

func Test_Pool_Start_HandlerPanicRecovered(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	handler.
		EXPECT().
		Handle(mock.AnythingOfType("*context.cancelCtx"), 1).
		RunAndReturn(func(_ context.Context, value int) int {
			panic(value)
		}).
		Once()

	collector := mockworkerpool.NewCollector[int](t)

	pool, err := workerpool.New(1 /*size*/, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()

	// act
	require.NoError(t, pool.Add(t.Context(), 1))

	// assert
	<-pool.Complete()
}

func Test_Pool_Start_CollectorPanicRecovered(t *testing.T) {
	// arrange
	handler := mockworkerpool.NewHandler[int, int](t)
	handler.
		EXPECT().
		Handle(mock.AnythingOfType("*context.cancelCtx"), mock.AnythingOfType("int")).
		RunAndReturn(func(_ context.Context, value int) int {
			return value * 2
		}).
		Once()

	collector := mockworkerpool.NewCollector[int](t)
	collector.
		EXPECT().
		Collect(2).
		RunAndReturn(func (value int) {
			panic(value)
		}).
		Once()

	pool, err := workerpool.New(1 /*size*/, handler, collector)
	require.NoError(t, err)

	go func() {
		pool.Start(t.Context())
	}()

	<-pool.Started()

	// act
	require.NoError(t, pool.Add(t.Context(), 1))

	// assert
	<-pool.Complete()
}
