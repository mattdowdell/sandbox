package main

import (
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/mattdowdell/sandbox/internal/drivers/term/components/tab"
	"github.com/mattdowdell/sandbox/internal/drivers/term/components/tabber"
)

//nolint:mnd // padding values
var (
	docStyle       = lipgloss.NewStyle().Padding(1, 2)
	highlightColor = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	bodyStyle      = lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop()
)

// ...
var _ tea.Model = (*Model)(nil)

// ...
type Model struct {
	tabber *tabber.Tabber
	height int
	width  int
}

// ...
func New() *Model {
	return &Model{
		tabber: tabber.New("Resources", "Events"),
	}
}

// ...
//
// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return tea.SetWindowTitle("Example Service Browser")
}

// ...
//
// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit

		case "right", "tab":
			m.tabber.MoveRight()

		case "left", "shift+tab":
			m.tabber.MoveLeft()
		}

	case tea.WindowSizeMsg:
		m.tabber.SetWindowSize(msg)

		m.height = msg.Height
		m.width = msg.Width
	}

	return m, nil
}

// ...
//
// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	doc := strings.Builder{}
	var rendered []string

	for i, t := range m.tabber.Tabs() {
		style := tab.Style(i == 0, t.Active, highlightColor)
		rendered = append(rendered, style.Render(t.Name))
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	gapLen := m.width - lipgloss.Width(row) - bodyStyle.GetHorizontalFrameSize() - docStyle.GetHorizontalFrameSize() - 1
	gap := tab.GapStyle(highlightColor).Render(strings.Repeat(" ", max(0, gapLen)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	doc.WriteString(row)
	doc.WriteString("\n")
	doc.WriteString(
		bodyStyle.
			Width(m.width - bodyStyle.GetHorizontalFrameSize() - docStyle.GetHorizontalFrameSize()).
			Height(m.height - bodyStyle.GetVerticalFrameSize() - docStyle.GetVerticalFrameSize() + 1).
			Render("content"),
	)

	return docStyle.Render(doc.String())
}
