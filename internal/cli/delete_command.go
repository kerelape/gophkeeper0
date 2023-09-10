package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type deleteCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*deleteCommand)(nil)

// Description implements command.
func (d *deleteCommand) Description() string {
	return "Delete resource."
}

// Help implements command.
func (d *deleteCommand) Help() string {
	return "<RID: int>"
}

// Execute implements command.
func (d *deleteCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("expected 1 arguments")
	}

	var identity, identityError = authenticate(ctx, d.gophkeeper)
	if identityError != nil {
		return true, identityError
	}

	var rid, ridError = strconv.Atoi(args.Pop())
	if ridError != nil {
		return false, ridError
	}
	if err := identity.Delete(ctx, (gophkeeper.ResourceID)(rid)); err != nil {
		return true, err
	}

	fmt.Printf("Successfully deleted resource (RID: %d),\n", rid)

	return true, nil
}
