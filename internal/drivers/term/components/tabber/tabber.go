package tabber

import (
	"github.com/charmbracelet/bubbletea"
)

// ...
type Tab struct {
	Name   string
	Active bool
}

// ...
type Tabber struct {
	tabs   []string
	active int
	height int
	width  int
}

// ...
func New(tabs ...string) *Tabber {
	return &Tabber{
		tabs: tabs,
	}
}

// ...
func (t *Tabber) SetWindowSize(msg tea.WindowSizeMsg) {
	t.width = msg.Width
	t.height = msg.Height
}

func (t *Tabber) MoveLeft() {
	t.active = (t.active - 1) % len(t.tabs)
}

func (t *Tabber) MoveRight() {
	t.active = (t.active + 1) % len(t.tabs)
}

func (t *Tabber) Tabs() []Tab {
	tabs := make([]Tab, 0, len(t.tabs))

	for i, tb := range t.tabs {
		tabs = append(tabs, Tab{
			Name:   tb,
			Active: i == t.active,
		})
	}

	return tabs
}
