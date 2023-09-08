package cli

import (
	"context"

	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
)

// CLI is gophkeeper command-line interface.
type CLI struct {
	Gophkeeper  gophkeeper.Gophkeeper
	CommandLine []string
}

var _ runnable.Runnable = (*CLI)(nil)

// Run implements runnable.Runnable.
//
// @todo #3 Implement `register` command.
// @todo #3 Implement `login` command.
// @todo #3 Implement `list` command.
// @todo #3 Implement `store-credential <description> <platform> <username> <password>` command.
// @todo #3 Implement `store-file <description> <path>` command.
// @todo #3 Implement `store-card <description> <cardholder> <number> <date> <cvv/cvc> command.
// @todo #3 Implement `store-text <description> <content>` command.
// @todo #3 Implement `restore-credential <rid>` command.
// @todo #3 Implement `restore-file <rid>` command.
// @todo #3 Implement `restore-card <rid>` command.
// @todo #3 Implement `restore-text <rid>` command.
func (c *CLI) Run(ctx context.Context) error {
	return nil
}
