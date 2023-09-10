package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type restoreTextCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*restoreTextCommand)(nil)

// Description implements command.
func (r *restoreTextCommand) Description() string {
	return "Restore text note."
}

// Help implements command.
func (r *restoreTextCommand) Help() string {
	return "<RID: int>"
}

// Execute implements command.
func (r *restoreTextCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("expected 1 argument")
	}

	var rid, ridError = strconv.Atoi(args.Pop())
	if ridError != nil {
		return false, ridError
	}

	var gophkeeperIdentity, gophkeeperIdentityError = authenticate(ctx, r.gophkeeper)
	if gophkeeperIdentityError != nil {
		return true, gophkeeperIdentityError
	}

	var vaultPassword, vaultPasswordError = vaultPassword(ctx)
	if vaultPasswordError != nil {
		return true, vaultPasswordError
	}

	var identity = identity{
		origin: gophkeeperIdentity,
	}
	var resource, resourceError = identity.RestoreText(ctx, (gophkeeper.ResourceID)(rid), vaultPassword)
	if resourceError != nil {
		return true, resourceError
	}

	fmt.Printf("\n%s\n", resource.content)

	return true, nil
}
