// Package repositories package describes the implementation-agnostic interaction points with other
// parts of the business. It formalises how the transfer of information works.
//
// While a repository can often represent a separate business system, it also encompasses data
// persistence. By ensuring both of these concerns are implementation-agnostic, it enables future
// replacement with minimal changes to business logic, or even selection between multiple options at
// runtime.
//
// A repository should aim to be a specific and targeted interface, rather than providing generic
// CRUDL operations. While the repository will almost certainly support CRUDL operations, a
// repisitory should be concerned with the business process rather than the application. For
// example, a repository may have a method for updating a specific column in a database row, despite
// multiple columns theoretically supporting updates.
//
// Some repositories might seem out of place initially, such as Clock and UUIDGenerator. These are
// more closely aligned with good application design principles rather than Clean Architecture. They
// are intended to abstract an external touch point of the application, namely interaction with the
// OS, to simplify testing and ensure portability. In this case they abstract getting the system
// time and reading random values respectively. This lack of following a rigid adherance is a good
// example of purism to the detriment of the application. Deviations are acceptable if they can be
// justified, and fundamental assumptions such as direction of dependencies are not broken.
package repositories
