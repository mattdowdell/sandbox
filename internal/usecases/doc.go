// Package usecases contains the usecases layer for the application.
//
// Usecases combine the business logic and models in the domain layer to provide a set of business
// processes for the application. They operate exclusively on the domain representation and are
// therefore agnostic to the input and output formats used by the application.
//
// The usecases layer codifies how the application performs the functions for which is has been
// designed and built. As a result, it can be used as a reference to understand what the application
// does. Furthermore, as each usecase is a separate file, a summary of application behaviour can be
// obtained by simply listing the files.
//
// # Stateless
//
// Each usecase is either stateless or contains state configured at runtime and not modified during
// the lifetime of the process. This allows a single instance of each usecase to be created at
// startup and injected into the rest of the application as needed.
//
// Any state that could change is provided when the usecase is executed. This includes
// request-scoped values such as database transactions. As a result, usecases are thread-safe and
// can be used concurrently without risk of race conditions.
//
// Deviations from this stateless approach can be made as long as the core assumption still fits.
// For example, a usecase field could internally manage its own runtime state as long as that state
// management is fully encapculated and thread-safe.
//
// # Error Boundary
//
// The usecases layer is a natural error boundary for the application. This means that errors that
// occur during the usecase execution are not bubbled up to caller of the usecase. Instead, a
// separate, clean error is returned that does not wrap any other error.
//
// This approach simplifies error handling in the caller as other packages do not need to be
// reviewed to identify the full range of error values. It also helps reinforce the idea that
// usecases are the entry and exit point of the business logic in the application. Finally, it
// means wrapped error types do not need to be considered when using [errors.Is] and [errors.As].
//
// # Consumer-Defined Interfaces
//
// In idiomatic Go, interfaces should be defined by the consumer and implicit interface
// implementation should be leveraged. The usecases here embrace that philosophy. This means that
// usecases avoid making any assumptions about how usecases are organised and invoked. It also
// allows for greater decoupling of the layers within the application.
package usecases
