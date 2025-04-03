package workerpool

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

// Errors that can be returned by the pool.
var (
	ErrInvalidSize = errors.New("size must be > 0")
	ErrQueueClosed = errors.New("worker pool queue has been closed")
	ErrFailedAdd   = errors.New("failed to add to queue")
	ErrFailedWait  = errors.New("failed to wait for workers")
)

// Handler implementations handle a work item from the queue. Implementations must be thread safe.
type Handler[T, U any] interface {
	Handle(context.Context, T) Result[U]
}

// Collector implementations receive results from the handler. Implementations do not need to be
// thread safe.
type Collector[U any] interface {
	Collect(Result[U])
}

// Pool provides a pool of workers. Each worker calls a handler to complete a work item. Results are
// then passed to a collector for aggregation.
//
// If Handler.Handle panics, the recovered value is converted to an error, wrapped in a Result, and
// passed to Collector.Collect as normal. The goroutine for the worker that saw the panic is then
// restarted.
//
// If Collector.Collect panics, the collector is restarted, but the collected value is lost.
type Pool[T, U any] struct {
	size      int
	handler   Handler[T, U]
	collector Collector[U]
	queue     chan T
	results   chan Result[U]
	started   chan struct{}
	waiting   chan struct{}
	complete  chan struct{}
	starter   sync.Once
	waiter    sync.Once
}

// New creates a new Pool.
func New[T, U any](
	size int,
	handler Handler[T, U],
	collector Collector[U],
) (*Pool[T, U], error) {
	if size <= 0 {
		return nil, ErrInvalidSize
	}

	return &Pool[T, U]{
		size:      size,
		handler:   handler,
		collector: collector,
		queue:     make(chan T),
		results:   make(chan Result[U]),
		started:   make(chan struct{}),
		waiting:   make(chan struct{}),
		complete:  make(chan struct{}),
	}, nil
}

// Start starts the workers in the pool. It blocks until all workers are stoped, either via Wait or
// Stop.
//
// This may be called multiple times without causing an error, but only the first call will start
// the workers and block.
func (p *Pool[T, U]) Start(ctx context.Context) {
	p.starter.Do(func() {
		go p.startCollector(ctx)

		var wg sync.WaitGroup

		for range p.size {
			wg.Add(1)
			go p.startWorker(ctx, &wg)
		}

		// signal that all threads have been started
		close(p.started)

		// wait for all threads to complete
		// then signal that the collector can stop
		wg.Wait()
		close(p.results)
	})
}

func (p *Pool[T, U]) startCollector(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "restarting panicked collector", slogx.Panic(r), slogx.Stacktrace())
			go p.startCollector(ctx)

			return
		}

		// signal that collection has completed
		close(p.complete)
	}()

	for r := range p.results {
		p.collector.Collect(r)
	}
}

func (p *Pool[T, U]) startWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		if r := recover(); r != nil {
			p.results <- Err[U](fmt.Errorf("recovered: %v", r))

			slog.ErrorContext(ctx, "restarting panicked worker", slogx.Panic(r), slogx.Stacktrace())
			go p.startWorker(ctx, wg)

			return
		}

		wg.Done()
	}()

	for {
		select {
		// handle the parent context being closed
		case <-ctx.Done():
			return

		case item, ok := <-p.queue:
			// handle Wait being called
			if !ok {
				return
			}

			result := p.handler.Handle(ctx, item)
			p.results <- result
		}
	}
}

// Add adds an item of work to the work queue and blocks until a worker thread has taken the item.
// ErrQueueClosed is returned if an item is added after Wait has been called.
func (p *Pool[T, U]) Add(ctx context.Context, item T) error {
	select {
	case <-p.waiting:
		return ErrQueueClosed

	default:
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrFailedAdd, ctx.Err())

		case p.queue <- item:
			return nil
		}
	}
}

// Wait waits until all the in-progress work is completed and should normally be called once all
// pending work items has been added to the queue. It can be called multiple times without error.
func (p *Pool[T, U]) Wait(ctx context.Context) error {
	p.waiter.Do(func() {
		// signal that no further work should be accepted
		close(p.waiting)

		// signal that no further work is incoming
		// then wait for the collector to process all results
		close(p.queue)
	})

	select {
	case <-p.complete:
		return nil

	case <-ctx.Done():
		return fmt.Errorf("%w: %w", ErrFailedWait, ctx.Err())
	}
}
