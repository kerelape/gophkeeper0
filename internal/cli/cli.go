package cli

import (
	"context"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kerelape/gophkeeper/internal/cli/application"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
)

// CLI is gophkeeper command-line interface.
type CLI struct {
	Server string // Gophkeeper REST server address.
}

var _ runnable.Runnable = (*CLI)(nil)

// Run implements runnable.Runnable.
func (c *CLI) Run(ctx context.Context) error {
	var program = tea.NewProgram(
		&application.Application{
			Gophkeeper: &gophkeeper.RestGophkeeper{
				Client: *http.DefaultClient,
				Server: c.Server,
			},
		},
	)
	go program.Run()
	<-ctx.Done()
	program.Quit()
	return nil
}
