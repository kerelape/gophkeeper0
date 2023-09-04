package cli

import (
	"context"

	"github.com/pior/runnable"
)

// CLI is gophkeeper command-line interface.
type CLI struct {
	Server string // Gophkeeper REST server address.
}

var _ runnable.Runnable = (*CLI)(nil)

// Run implements runnable.Runnable.
//
// @todo #3 Implement CLI.
func (c *CLI) Run(ctx context.Context) error {
	panic("unimplemented")
}
