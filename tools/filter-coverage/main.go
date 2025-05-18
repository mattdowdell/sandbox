package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"strings"
	"os"
	"go/token"
	"go/ast"
	"go/parser"
	"path/filepath"

	"github.com/mattdowdell/sandbox/pkg/slogx"
)

const (
	exitSuccess = iota
	exitFailure
)

var (
	module = flag.String("mod", "go.mod", "The path of the go.mod file for the module under test")
	output = flag.String("output", "filtered.out", "The path to output to")
)

func main() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <coverprofile>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		slog.ErrorContext(ctx, "missing argument")
		return exitFailure
	}

	profile, err := os.Open(flag.Arg(0))
	if err != nil {
		slog.ErrorContext(ctx, "failed to open coverprofile", slogx.Err(err))
		return exitFailure
	}
	defer profile.Close()

	modPath, err := extractModPath(*module)
	if err != nil {
		slog.ErrorContext(ctx, "failed to extract module path", slogx.Err(err))
		return exitFailure
	}

	scanner := bufio.NewScanner(profile)
	processed := map[string]bool{}
	fset := token.NewFileSet()

	var buffer bytes.Buffer

	for scanner.Scan() {
		line := scanner.Text()

		filename, ok := extractFilename(line, modPath)
		if !ok {
			buffer.WriteString(line + "\n")
			continue
		}

		generated, ok := processed[filename]
		if !ok {
			generated = isGenerated(fset, filename)
			processed[filename] = generated
		}

		if !generated {
			buffer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		slog.ErrorContext(ctx, "failed to read file", slogx.Err(err))
		return exitFailure
	}

	out, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		slog.ErrorContext(ctx, "failed to open output file", slogx.Err(err))
		return exitFailure
	}
	defer out.Close()

	if _, err := out.WriteString(buffer.String()); err != nil {
		slog.ErrorContext(ctx, "failed to write to output file", slogx.Err(err))
		return exitFailure
	}

	return exitSuccess
}

func extractModPath(path string) (string, error) {
	return "github.com/mattdowdell/sandbox/", nil
}

func extractFilename(line, prefix string) (string, bool) {
	before, _, ok := strings.Cut(line, ":")
	if !ok {
		return "", false
	}

	filename, _ := strings.CutPrefix(before, prefix)

	return filename, true
}

func isGenerated(fset *token.FileSet, filename string) bool {
	parsed, err := parser.ParseFile(fset, filename, nil /*src*/, parser.ParseComments /*mode*/)
	if err != nil {
		return false
	}

	return ast.IsGenerated(parsed)
}
