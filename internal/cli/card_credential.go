package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type cardInfo struct {
	ccn    string
	exp    string
	cvv    string
	holder string
}

func cardCredential(ctx context.Context) (cardInfo, error) {
	var m, err = tea.NewProgram(
		newCardCredentialModel(),
		tea.WithAltScreen(),
		tea.WithContext(ctx),
	).Run()
	if err != nil {
		return cardInfo{}, err
	}
	var cm = m.(cardCredentialModel)
	var info = cardInfo{
		ccn:    cm.ccn.Value(),
		exp:    cm.exp.Value(),
		cvv:    cm.cvv.Value(),
		holder: cm.holder.Value(),
	}
	return info, nil
}

type cardCredentialModel struct {
	width, height int

	ccn    textinput.Model
	exp    textinput.Model
	cvv    textinput.Model
	holder textinput.Model

	cancelled bool
}

func newCardCredentialModel() cardCredentialModel {
	var m = cardCredentialModel{
		ccn:    textinput.New(),
		exp:    textinput.New(),
		cvv:    textinput.New(),
		holder: textinput.New(),
	}

	m.ccn.CharLimit = 16 + 3
	m.ccn.Placeholder = "4505 **** **** 1234"
	m.ccn.Prompt = ""
	m.ccn.Validate = func(s string) error {
		if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
			return fmt.Errorf("CCN is invalid")
		}
		if len(s)%5 == 0 && s[len(s)-1] != ' ' {
			return fmt.Errorf("CCN must separate groups with spaces")
		}
		c := strings.ReplaceAll(s, " ", "")
		_, err := strconv.ParseInt(c, 10, 64)
		return err
	}
	m.ccn.Focus()

	m.exp.CharLimit = 5
	m.exp.Placeholder = "MM/YY"
	m.exp.Prompt = ""
	m.exp.Validate = func(s string) error {
		e := strings.ReplaceAll(s, "/", "")
		_, err := strconv.ParseInt(e, 10, 64)
		if err != nil {
			return fmt.Errorf("EXP is invalid")
		}
		if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
			return fmt.Errorf("EXP is invalid")
		}
		return nil
	}

	m.cvv.CharLimit = 3
	m.cvv.EchoMode = textinput.EchoPassword
	m.cvv.Placeholder = "123"
	m.cvv.Prompt = ""
	m.cvv.Validate = func(s string) error {
		_, err := strconv.ParseInt(s, 10, 64)
		return err
	}

	m.holder.CharLimit = 64
	m.holder.Placeholder = "CARD HOLDER"
	m.holder.Prompt = ""

	return m
}

var _ tea.Model = (*cardCredentialModel)(nil)

// Init implements tea.Model.
func (m cardCredentialModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m cardCredentialModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cnnCmd    tea.Cmd
		expCmd    tea.Cmd
		cvvCmd    tea.Cmd
		holderCmd tea.Cmd
	)
	m.ccn, cnnCmd = m.ccn.Update(msg)
	m.exp, expCmd = m.exp.Update(msg)
	m.cvv, cvvCmd = m.cvv.Update(msg)
	m.holder, holderCmd = m.holder.Update(msg)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch {
			case m.ccn.Focused():
				if len(m.ccn.Value()) > 0 && m.ccn.Err == nil {
					m.ccn.Blur()
					m.exp.Focus()
				}
			case m.exp.Focused():
				if len(m.exp.Value()) > 0 && m.exp.Err == nil {
					m.exp.Blur()
					m.cvv.Focus()
				}
			case m.cvv.Focused():
				if len(m.cvv.Value()) > 0 && m.cvv.Err == nil {
					m.cvv.Blur()
					m.holder.Focus()
				}
			case m.holder.Focused():
				m.holder.Blur()
				return m, tea.Quit
			}
		case "esc", "ctrl+c":
			m.cancelled = true
			return m, tea.Quit
		}
	}
	return m, tea.Batch(cnnCmd, expCmd, cvvCmd, holderCmd)
}

// View implements tea.Model.
func (m cardCredentialModel) View() string {
	var hintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#002288"))
	return form(
		m.width, m.height,
		"Card",
		lipgloss.JoinVertical(
			lipgloss.Left,
			hintStyle.Render("Card Number"),
			m.ccn.View(),
			" ",
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				lipgloss.JoinVertical(
					lipgloss.Left,
					hintStyle.Render("Expiration Date"),
					m.exp.View(),
				),
				"       ",
				lipgloss.JoinVertical(
					lipgloss.Left,
					hintStyle.Render("CVV"),
					m.cvv.View(),
				),
			),
			" ",
			hintStyle.Render("Card Holder"),
			m.holder.View(),
			" ",
			help.New().ShortHelpView(
				[]key.Binding{
					key.NewBinding(
						key.WithKeys("esc"),
						key.WithHelp("[esc]", "cancel"),
					),
					key.NewBinding(
						key.WithKeys("enter"),
						key.WithHelp("[enter]", "submit"),
					),
				},
			),
		),
	)
}
