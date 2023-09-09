package cli

import "github.com/charmbracelet/lipgloss"

var formStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("#0088AA")).
	Padding(1, 2, 1, 2)

func form(width, height int, title, content string) string {
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			title,
			formStyle.Render(content),
		),
	)
}
