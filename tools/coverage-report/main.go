package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/mattdowdell/sandbox/pkg/exit"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

var (
	flagCurrent  = flag.String("current", "", "The current coverage results file.")
	flagPrevious = flag.String("previous", "", "The previous coverage results file, if available.")
)

func main() {
	os.Exit(run())
}

func run() int {
	setupDefaultLogger()

	flag.Parse()

	current, err := PackageCoverageFromFile(*flagCurrent)
	if err != nil {
		slog.Error("failed to parse current profiles", slogx.Err(err))
		return exit.Failure
	}

	previous, err := PackageCoverageFromFile(*flagPrevious)
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
			change.GetPrevious(),
			change.GetCurrent(),
			change.Diff(),
		)
	}

	return exit.Success
}

func setupDefaultLogger() {
	logger := slog.New(slog.NewTextHandler(flag.CommandLine.Output(), &slog.HandlerOptions{
		Level: slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) > 0 {
				return a
			}

			if a.Key == slog.TimeKey && a.Value.Kind() == slog.KindTime {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.Kitchen))
			}

			return a
		},
	}))

	slog.SetDefault(logger)
}
