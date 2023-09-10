package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type textModel struct {
	width, height int

	content   textarea.Model
	cancelled bool

	help help.Model
}

func newTextModel() textModel {
	var m = textModel{
		content: textarea.New(),
		help:    help.New(),
	}
	m.content.ShowLineNumbers = false
	m.content.MaxWidth = 64
	m.content.CharLimit = 1024
	m.content.Placeholder = "type your note..."
	m.content.SetHeight(12)
	m.content.SetWidth(64)
	m.content.Focus()
	return m
}

func text(ctx context.Context) (string, error) {
	var m, err = tea.NewProgram(
		newTextModel(),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return "", err
	}
	if m.(textModel).cancelled {
		return "", errors.New("user cancelled typing")
	}
	return m.(textModel).content.Value(), nil
}

var _ tea.Model = (*textModel)(nil)

// Init implements tea.Model.
func (m textModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update implements tea.Model.
func (m textModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var descriptionCmd tea.Cmd
	m.content, descriptionCmd = m.content.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
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
func (m textModel) View() string {
	return form(
		m.width, m.height,
		"Text Note",
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.content.View(),
			lipgloss.PlaceHorizontal(
				lipgloss.Width(m.content.View()),
				lipgloss.Right,
				fmt.Sprintf("%d/%d", m.content.Length(), m.content.CharLimit),
			),
			" ",
			m.help.ShortHelpView(
				[]key.Binding{
					key.NewBinding(
						key.WithKeys("esc"),
						key.WithHelp("[esc]", "cancel"),
					),
					key.NewBinding(
						key.WithKeys("ctrl+s"),
						key.WithHelp("[ctrl+s]", "save"),
					),
				},
			),
		),
	)
}
