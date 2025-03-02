package repositories

// ...
type Pager struct {
	Limit int
}

// ...
type Paged[T any] struct {
	Items []T
}
