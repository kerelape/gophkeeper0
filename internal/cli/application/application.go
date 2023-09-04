package application

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Application is gophkeeper CLI application.
//
// @todo #54 Implement Application.
type Application struct {
	Gophkeeper gophkeeper.Gophkeeper
}

var _ tea.Model = (*Application)(nil)

// Init implements tea.Model.
func (m *Application) Init() tea.Cmd {
	panic("unimplemented")
}

// Update implements tea.Model.
func (m *Application) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	panic("unimplemented")
}

// View implements tea.Model.
func (m *Application) View() string {
	panic("unimplemented")
}
