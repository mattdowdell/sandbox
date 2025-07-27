package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// Change represents the coverage change between 2 coverage results.
type Change struct {
	Name     string
	Previous *float64
	Current  *float64
}

// CollectChanges combines 2 coverage results into a series of per-package changes.
func CollectChanges(current, previous map[string]*Package) []*Change {
	merged := slices.Concat(
		slices.Collect(maps.Keys(current)),
		slices.Collect(maps.Keys(previous)),
	)

	slices.Sort(merged)
	names := slices.Compact(merged)

	changes := make([]*Change, 0, len(names))

	for _, name := range names {
		changes = append(changes, &Change{
			Name:     name,
			Previous: previous[name].CoveragePercent(),
			Current:  current[name].CoveragePercent(),
		})
	}

	slices.SortFunc(changes, func(a, b *Change) int {
		if a.Name == totalName {
			return 1
		}

		if b.Name == totalName {
			return -1
		}

		return cmp.Compare(a.Name, b.Name)
	})

	return changes
}

// GetPrevious returns the previous coverage for the package formatted as a percentage or "-" if no
// value was found.
func (c *Change) GetPrevious() string {
	if c.Previous == nil {
		return "-"
	}

	return fmt.Sprintf("%.1f%%", *c.Previous)
}

// GetCurrent returns the current coverage for the package formatted as a percentage or "-" if no
// value was found.
func (c *Change) GetCurrent() string {
	if c.Current == nil {
		return "-"
	}

	return fmt.Sprintf("%.1f%%", *c.Current)
}

// Diff returns the coverage change as a percentage. If either the current or previous coverage were
// not found, "N/A" is returned.
func (c *Change) Diff() string {
	if c.Previous == nil || c.Current == nil {
		return "N/A"
	}

	return fmt.Sprintf("%+.1f%%", *c.Current-*c.Previous)
}
