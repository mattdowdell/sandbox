package main

import (
	"fmt"
	"path"

	"golang.org/x/tools/cover"
)

const (
	totalName = "Total"
)

// Package represents the number of statements that have test coverage within a single package.
type Package struct {
	Name    string
	Total   int
	Covered int
}

// PackageCoverageFromFile parses a file containing per-file coverage profiles and converts it into
// per-package representation.
func PackageCoverageFromFile(filename string) (map[string]*Package, error) {
	if filename == "" {
		return nil, nil
	}

	profiles, err := cover.ParseProfiles(*flagCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profiles: %w", err)
	}

	return PackageCoverage(profiles), err
}

// PackageCoverage converts per-file coverage profiles into per-package representation.
func PackageCoverage(profiles []*cover.Profile) map[string]*Package {
	packages := map[string]*Package{}

	total := &Package{}

	for _, profile := range profiles {
		name := path.Dir(profile.FileName)

		pkg, ok := packages[name]
		if !ok {
			pkg = &Package{
				Name: name,
			}
		}

		if pkg.Name == "" {
			pkg.Name = name
		}

		for _, block := range profile.Blocks {
			pkg.Total += block.NumStmt
			total.Total += block.NumStmt

			if block.Count > 0 {
				pkg.Covered += block.NumStmt
				total.Covered += block.NumStmt
			}
		}

		packages[name] = pkg
	}

	packages[totalName] = total

	return packages
}

// CoveragePercent returns the statement coverage for the package as a percentage.
//
//nolint:mnd // calculates a percentage
func (p *Package) CoveragePercent() *float64 {
	if p == nil {
		return nil
	}

	value := (float64(p.Covered) / float64(p.Total)) * 100.0
	return &value
}
