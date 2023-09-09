package cli

import (
	"context"
	"errors"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func vaultPassword(ctx context.Context) (string, error) {
	var m, err = tea.NewProgram(
		newVaultPasswordModel(),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return "", err
	}
	if m.(vaultPasswordModel).cancelled {
		return "", errors.New("vault password typing cancelled by user")
	}
	return m.(vaultPasswordModel).password.Value(), nil
}

type vaultPasswordModel struct {
	width, height int
	cancelled     bool

	password textinput.Model
}

func newVaultPasswordModel() vaultPasswordModel {
	var m = vaultPasswordModel{
		password: textinput.New(),
	}
	m.password.EchoMode = textinput.EchoPassword
	m.password.Prompt = "Vault password: "
	m.password.CharLimit = 32
	m.password.Placeholder = "enter your vault password..."
	m.password.Focus()
	return m
}

var _ tea.Model = (*vaultPasswordModel)(nil)

// Init implements tea.Model.
func (v vaultPasswordModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (v vaultPasswordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var passwordCmd tea.Cmd
	v.password, passwordCmd = v.password.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		v.height = msg.Height
		v.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			v.password.Blur()
			return v, tea.Quit
		case "ctrl+c", "esc":
			v.cancelled = true
			return v, tea.Quit
		}
	}
	return v, tea.Batch(passwordCmd)
}

// View implements tea.Model.
func (v vaultPasswordModel) View() string {
	return form(
		v.width, v.height,
		"Vault",
		lipgloss.NewStyle().Width(64).Render(v.password.View()),
	)
}
