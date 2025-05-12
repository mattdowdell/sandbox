// Package domain contains the domain layer for the application.
//
// The domain layer models the business dependent but application independent code. It contains the
// representation of the real world objects and processes that the application intends to model.
// This layer is also the centre of the application, meaning that it should not include any
// dependencies on the other three layers of the application.
//
// The domain layer is usually quite thin, and sometimes only contains some struct and repository
// definitions. However, representations of business processes also live in this layer.
//
// When deciding if code should go in the domain layer, a good starting question is "does this code
// model something that would exist in the business even if our application did not exist?". If the
// answer is yes, then it most likely lives here.
package domain
