package workerpool

// Result contains a success or error returned by a Handler.
type Result[U any] struct {
	OK  U
	Err error
}

// OK creates a success Result.
func OK[U any](value U) Result[U] {
	return Result[U]{
		OK: value,
	}
}

// Err creates an error Result.
func Err[U any](err error) Result[U] {
	return Result[U]{
		Err: err,
	}
}

// IsOK returns true if the Err field of the result is nil.
func (r *Result[U]) IsOK() bool {
	return r.Err == nil
}

// IsErr returns true if the Err field of the result is not nil.
func (r *Result[U]) IsErr() bool {
	return r.Err != nil
}
