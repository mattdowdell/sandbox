package main

import (
	"github.com/charmbracelet/bubbletea"
)

// ...
var _ tea.Model = (*Model)(nil)

// ...
type Model struct {}

// ...
func New() *Model {
	return &Model{}
}

// ...
//
// Init is the first function that will be called. It returns an optional
// initial command. To not perform an initial command return nil.
func (m *Model) Init() tea.Cmd {
	return nil
}

// ...
//
// Update is called when a message is received. Use it to inspect messages
// and, in response, update the model and/or send a command.
func (m *Model) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return nil, nil // TODO: panics here?
}

// ...
//
// View renders the program's UI, which is just a string. The view is
// rendered after every Update.
func (m *Model) View() string {
	return "not implemented"
}
