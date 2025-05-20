package main

import (
	"os"
)

const (
	exitSuccess = iota
	exitFailure
)

func main() {
	os.Exit(run())
}

func run() int {
	return exitSuccess
}
