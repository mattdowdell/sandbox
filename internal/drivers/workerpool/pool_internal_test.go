package workerpool

// Started returns a channel that is closed when all worker threads have started.
//
// This is intended for use in testing to reliably determine the effects of Wait.
func (p *Pool[T, U]) Started() <-chan struct{} {
	return p.started
}

// Started returns a channel that is closed when all worker threads have stopped.
//
// This is intended for use in testing so tests can reliably determine the effects of stopping or
// waiting.
func (p *Pool[T, U]) Complete() <-chan struct{} {
	return p.complete
}
