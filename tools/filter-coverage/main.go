package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/mod/modfile"

	"github.com/mattdowdell/sandbox/internal/drivers/exit"
	"github.com/mattdowdell/sandbox/pkg/slogx"
)

var (
	modfilePath = flag.String("modfile", "go.mod", "The path of the go.mod file for the module under test")
	output      = flag.String("output", "filtered.out", "The path to output to")
	debug       = flag.Bool("debug", false, "Enable debug logging")
)

var level slog.LevelVar

func main() {
	os.Exit(run())
}

func run() int {
	setupDefaultLogger()

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <coverprofile>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if *debug {
		level.Set(slog.LevelDebug)
	}

	if flag.NArg() < 1 {
		slog.Error("missing argument")
		return exit.Failure
	}

	profile, err := os.Open(flag.Arg(0))
	if err != nil {
		slog.Error("failed to open coverprofile", slogx.Err(err))
		return exit.Failure
	}
	defer profile.Close()

	prefix, err := extractModulePrefix(*modfilePath)
	if err != nil {
		slog.Error("failed to extract module path", slogx.Err(err))
		return exit.Failure
	}

	filtered, err := filterFile(profile, prefix)
	if err != nil {
		slog.Error("failed to filter file", slogx.Err(err))
		return exit.Failure
	}

	out, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("failed to open output file", slogx.Err(err))
		return exit.Failure
	}
	defer out.Close()

	if _, err := out.WriteString(filtered); err != nil {
		slog.Error("failed to write to output file", slogx.Err(err))
		return exit.Failure
	}

	return exit.Success
}

func setupDefaultLogger() {
	logger := slog.New(slog.NewTextHandler(flag.CommandLine.Output(), &slog.HandlerOptions{
		Level: &level,
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

func filterFile(profile *os.File, modPath string) (string, error) {
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

			if generated {
				slog.Debug("skipping generated file", slog.String("filename", filename))
			}
		}

		if !generated {
			buffer.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func extractModulePrefix(path string) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	if name := modfile.ModulePath(contents); name != "" {
		return name + "/", nil
	}

	return "", fmt.Errorf("failed to extract module name from: %s", path)
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
