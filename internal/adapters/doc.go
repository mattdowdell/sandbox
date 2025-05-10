// Package adapters provides application-specific logic.
//
// Adapters contain code that exists for the application only. It is responsible for converting data
// in external formats to and from the internal domain representation. Adapters will often
// correspond to an application entrypoint or domain repository implementation, although this is not
// always the case.
//
// Adapters that provide application entrypoints are responsible for executing usecases. A usecase
// can be re-used by multiple adapters if overlap is required. This is useful if there are APIs in
// different format, such as gRPC and REST. Similarly, it can be useful if supporting multiple
// versions of the same API where each version has small differences in the external format.
//
// An adapter may build upon another adapter that also provides application dependent, but business
// independent code. Examples of this can be found in the common and usecasefacades adapter
// packages.
package adapters
