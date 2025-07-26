package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

var (
	flagCurrent  = flag.String("current", "", "TODO")
	flagPrevious = flag.String("previous", "", "TODO")
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()

	current, err := GatherCoverageFromFile(*flagCurrent)
	if err != nil {
		slog.Error("failed to parse current profiles", slogx.Err(err))
		return exit.Failure
	}

	previous, err := GatherCoverageFromFile(*flagPrevious)
	if err != nil {
		slog.Error("failed to parse previous profiles", slogx.Err(err))
		return exit.Failure
	}

	changes := CollectChanges(current, previous)

	fmt.Println("| Name | Previous | Current | Diff |")
	fmt.Println("| ---- | -------- | ------- | ---- |")

	for _, change := range changes {
		fmt.Printf(
			"| %s | %s | %s | %s |\n",
			change.Name,
			change.GetBefore(),
			change.GetAfter(),
			change.Diff(),
		)
	}

	return exit.Success
}
