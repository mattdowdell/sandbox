package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/internal/drivers/logging"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	logger := logging.NewAsDefault(slog.LevelInfo)

	program := tea.NewProgram(New(), tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		logger.ErrorContext(ctx, "failed to start program", slogx.Err(err))
		return exit.Failure
	}

	return exit.Success
}
