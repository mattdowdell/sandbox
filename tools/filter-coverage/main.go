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

	"golang.org/x/mod/modfile"
)

const (
	exitSuccess = iota
	exitFailure
)

var (
	modfilePath = flag.String("modfile", "go.mod", "The path of the go.mod file for the module under test")
	output      = flag.String("output", "filtered.out", "The path to output to")
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s <coverprofile>\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		slog.Error("missing argument")
		return exitFailure
	}

	profile, err := os.Open(flag.Arg(0))
	if err != nil {
		slog.Error("failed to open coverprofile", slog.Any("err", err))
		return exitFailure
	}
	defer profile.Close()

	prefix, err := extractModulePrefix(*modfilePath)
	if err != nil {
		slog.Error("failed to extract module path", slog.Any("err", err))
		return exitFailure
	}

	filtered, err := filterFile(profile, prefix)
	if err != nil {
		slog.Error("failed to filter file", slog.Any("err", err))
		return exitFailure
	}

	out, err := os.OpenFile(*output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("failed to open output file", slog.Any("err", err))
		return exitFailure
	}
	defer out.Close()

	if _, err := out.WriteString(filtered); err != nil {
		slog.Error("failed to write to output file", slog.Any("err", err))
		return exitFailure
	}

	return exitSuccess
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
				slog.Info("skipping generated file", slog.String("filename", filename))
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
