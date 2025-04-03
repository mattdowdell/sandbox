// Package workerpool provides a thread pool of workers.
//
// The worker pool creates a number of threads to process work items based on a pre-determined size.
// Once the workers have started, the work items can be added to the queue. Each worker receives an
// item from the queue and then passes it to a handler for processing. The handler returns the
// result of the processing, which is then passed to the collector to aggregate the results.
//
// Once all work items have been added to the queue, the pool can be waited upon until the collector
// has been provided all the results.
package workerpool
