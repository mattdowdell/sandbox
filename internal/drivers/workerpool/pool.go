// Package workerpool provides an experimental thread pool of workers.
//
// The worker pool creates a number of threads to process work items based on a pre-determined size.
// Once the workers have started, the work items can be added to the queue. Each worker receives an
// item from the queue and then passes it to a handler for processing. The handler returns the
// result of the processing, which is then passed to the collector to aggregate the results.
//
// Once all work items have been added to the queue, the pool can be waited upon until the collector
// has been provided all the results.
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
	ErrFailedWait  = errors.New("failed to wait for workers")
)

// Handler implementations handle a work item from the queue. Implementations must be thread safe.
type Handler[T, U any] interface {
	Handle(context.Context, T) U
}

// Collector implementations receive results from the handler. Implementations do not need to be
// thread safe.
type Collector[U any] interface {
	Collect(U)
}

// Pool provides a pool of workers. Each worker calls a handler to complete a work item. Results are
// then passed to a collector for aggregation.
type Pool[T, U any] struct {
	size      int
	handler   Handler[T, U]
	collector Collector[U]
	queue     chan T
	results   chan U
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
		results:   make(chan U),
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
			slog.ErrorContext(ctx, "collector panicked, restarting", slogx.Panic(r), slogx.Stacktrace())

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
	fmt.Println("startWorker")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "worker panicked, restarting", slogx.Panic(r)) // , slogx.Stacktrace())

			wg.Add(1)
			go p.startWorker(ctx, wg)
		}

		wg.Done()
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		fmt.Printf("working: %#v\n", p.queue)
		select {
		// handle the parent context being closed
		case <-ctx.Done():
			return

		case item, ok := <-p.queue:
			fmt.Println("item:", item, ok)
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
func (p *Pool[T, U]) Add(item T) error {
	fmt.Println("adding", item)
	select {
	case <-p.waiting:
		return ErrQueueClosed

	default:
		p.queue <- item
		fmt.Println("added", item)
		return nil
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
