package tab

import (
	"github.com/charmbracelet/lipgloss"
)

// ...
func GapStyle(color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(gapBorder(), true /*sides*/).
		BorderLeft(false).
		BorderTop(false).
		BorderForeground(color).
		Padding(0, 1)
}

// ...
func Style(first, active bool, color lipgloss.TerminalColor) lipgloss.Style {
	if active {
		return activeStyle(first, color)
	}

	return inactiveStyle(first, color)
}

// ...
func activeStyle(first bool, color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Border(activeBorder(first), true /*sides*/).
		BorderForeground(color).
		Padding(0, 1)
}

// ...
func inactiveStyle(first bool, color lipgloss.TerminalColor) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(inactiveBorder(first), true /*sides*/).
		BorderForeground(color).
		Padding(0, 1)
}

// ...
func activeBorder(first bool) lipgloss.Border {
	b := lipgloss.RoundedBorder()

	if first {
		b.BottomLeft = "│"
	} else {
		b.BottomLeft = "┘"
	}

	b.Bottom = " "
	b.BottomRight = "└"

	return b
}

// ...
func inactiveBorder(first bool) lipgloss.Border {
	b := lipgloss.RoundedBorder()

	if first {
		b.BottomLeft = "├"
	} else {
		b.BottomLeft = "┴"
	}

	b.BottomRight = "┴"

	return b
}

// ...
func gapBorder() lipgloss.Border {
	b := lipgloss.RoundedBorder()

	b.BottomRight = "┐"
	b.Right = " "

	return b
}
