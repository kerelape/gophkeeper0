package cli

import (
	"context"
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type credentialModel struct {
	width, height int
	cancelled     bool

	username textinput.Model
	password textinput.Model
}

func newCredentialModel() credentialModel {
	var m = credentialModel{
		username: textinput.New(),
		password: textinput.New(),
	}
	m.username.CharLimit = 32
	m.username.Prompt = "Username: "
	m.username.Placeholder = "type username..."
	m.username.Focus()

	m.password.CharLimit = 32
	m.password.Prompt = "Password: "
	m.password.Placeholder = "type password..."
	m.password.EchoMode = textinput.EchoPassword
	return m
}

func credential(ctx context.Context) (string, string, error) {
	var m, err = tea.NewProgram(
		newCredentialModel(),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return "", "", err
	}
	if m.(credentialModel).cancelled {
		return "", "", errors.New("credential form cancelled by user")
	}
	var credential = m.(credentialModel)
	return credential.username.Value(), credential.password.Value(), nil
}

var _ tea.Model = (*credentialModel)(nil)

// Init implements tea.Model.
func (m credentialModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m credentialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		usernameCmd tea.Cmd
		passwordCmd tea.Cmd
	)
	m.username, usernameCmd = m.username.Update(msg)
	m.password, passwordCmd = m.password.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch {
			case m.username.Focused():
				if m.username.Value() == "" {
					return m, textinput.Blink
				}
				m.username.Blur()
				m.password.Focus()
			case m.password.Focused():
				if m.password.Value() == "" {
					return m, textinput.Blink
				}
				m.password.Blur()
				return m, tea.Quit
			}
		case "tab":
			if m.password.EchoMode == textinput.EchoPassword {
				m.password.EchoMode = textinput.EchoNormal
			} else {
				m.password.EchoMode = textinput.EchoPassword
			}
		case "ctrl+c", "esc":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, tea.Batch(usernameCmd, passwordCmd)
}

// View implements tea.Model.
func (m credentialModel) View() string {
	var help = help.New()
	help.Width = 64
	return form(
		m.width, m.height,
		"Credential",
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.username.View(),
			strings.Repeat(" ", 64),
			m.password.View(),
			strings.Repeat(" ", 64),
			help.ShortHelpView(
				[]key.Binding{
					key.NewBinding(
						key.WithHelp("[esc]", "cancel"),
						key.WithKeys("esc"),
					),
					key.NewBinding(
						key.WithHelp("[tab]", "view password"),
						key.WithKeys("tab"),
					),
				},
			),
		),
	)
}
