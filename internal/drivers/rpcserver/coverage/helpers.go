package coverage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/coverage"
)

func writeCoverage() error {
	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		return errors.New("GOCOVERDIR not set")
	}

	if err := coverage.WriteCountersDir(dir); err != nil {
		return fmt.Errorf("failed to write covcounters: %w", err)
	}

	if err := coverage.WriteMetaDir(dir); err != nil {
		return fmt.Errorf("failed to write covmeta: %w", err)
	}

	if err := coverage.ClearCounters(); err != nil {
		return fmt.Errorf("failed to clear coverage counters: %w", err)
	}

	return nil
}

func listCoverage() ([]string, error) {
	dir := os.Getenv("GOCOVERDIR")
	if dir == "" {
		return nil, errors.New("GOCOVERDIR not set")
	}

	return filepath.Glob(filepath.Join(dir, "*"))
}
