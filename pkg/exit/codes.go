// Package exit provides constants for process exit codes.
package exit

// Exit codes to pass to [os.Exit] to exit a process.
//
// The values here are deliberately limited as they are intended for use in microservices and exit
// codes are rarely useful beyoind a simple success or failure. For a richer set of values, consider
// adopting those defined in [sysexits(3)] when interactively communicating runtime errors.
//
// [os.Exit]: https://pkg.go.dev/os#Exit
// [sysexits(3)]: https://man.freebsd.org/cgi/man.cgi?query=sysexits
const (
	Success = iota
	Failure
)
