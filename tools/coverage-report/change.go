package main

import (
	"cmp"
	"fmt"
	"maps"
	"slices"
)

// ...
type Change struct {
	Name   string
	Before *float64
	After  *float64
}

// ...
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
			Name:   name,
			Before: previous[name].Coverage(),
			After:  current[name].Coverage(),
		})
	}

	slices.SortFunc(changes, func(a, b *Change) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return changes
}

// ...
func (c *Change) GetBefore() string {
	if c.Before == nil {
		return "-"
	}

	return fmt.Sprintf("%.1f%%", *c.Before)
}

// ...
func (c *Change) GetAfter() string {
	if c.After == nil {
		return "-"
	}

	return fmt.Sprintf("%.1f%%", *c.After)
}

// ...
func (c *Change) Diff() string {
	if c.Before == nil || c.After == nil {
		return "N/A"
	}

	return fmt.Sprintf("%.1f%%", *c.After-*c.Before)
}

func (c *Change) String() string {
	return fmt.Sprintf("%s: %s -> %s (%s)", c.Name, c.GetBefore(), c.GetAfter(), c.Diff())
}
