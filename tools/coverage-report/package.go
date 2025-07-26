package main

import (
	"fmt"
	"path"

	"golang.org/x/tools/cover"
)

// ...
type Package struct {
	Name    string
	Total   int
	Covered int
}

// ...
func GatherCoverageFromFile(filename string) (map[string]*Package, error) {
	if filename == "" {
		return nil, nil
	}

	profiles, err := cover.ParseProfiles(*flagCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse profiles: %w", err)
	}

	return GatherCoverage(profiles), err
}

// ...
func GatherCoverage(profiles []*cover.Profile) map[string]*Package {
	packages := map[string]*Package{}

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
			if block.Count > 0 {
				pkg.Covered += block.NumStmt
			}
		}

		packages[name] = pkg
	}

	return packages
}

// ...
//
//nolint:mnd // calculates a percentage
func (p *Package) Coverage() *float64 {
	if p == nil {
		return nil
	}

	value := (float64(p.Covered) / float64(p.Total)) * 100.0
	return &value
}
