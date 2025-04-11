package repositories

// Pager contains the inputs for pagination.
type Pager struct {
	Limit int
}

// Paged contains a single page of results.
type Paged[T any] struct {
	Items []T
}
