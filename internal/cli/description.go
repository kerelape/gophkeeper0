package cli

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type descriptionModel struct {
	width, height int

	description textarea.Model
	cancelled   bool
}

func newDescriptionModel() descriptionModel {
	var m = descriptionModel{
		description: textarea.New(),
	}
	m.description.ShowLineNumbers = false
	m.description.MaxHeight = 8
	m.description.MaxWidth = 32
	m.description.Placeholder = "type description..."
	m.description.SetHeight(8)
	m.description.SetWidth(32)
	m.description.Focus()
	return m
}

func description(ctx context.Context) (string, error) {
	var m, err = tea.NewProgram(
		newDescriptionModel(),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return "", err
	}
	if m.(descriptionModel).cancelled {
		return "", errors.New("user cancelled typing description")
	}
	return m.(descriptionModel).description.Value(), nil
}

var _ tea.Model = (*descriptionModel)(nil)

// Init implements tea.Model.
func (m descriptionModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update implements tea.Model.
func (m descriptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var descriptionCmd tea.Cmd
	m.description, descriptionCmd = m.description.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			return m, tea.Quit
		case "ctrl+c", "esc":
			m.cancelled = false
			return m, tea.Quit
		}
	}
	return m, tea.Batch(descriptionCmd)
}

// View implements tea.Model.
func (m descriptionModel) View() string {
	return form(
		m.width, m.height,
		"Description",
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.description.View(),
			"ctrl+s to submit",
		),
	)
}
