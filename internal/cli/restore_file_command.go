package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type restoreFileCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*restoreFileCommand)(nil)

// Description implements command.
func (r *restoreFileCommand) Description() string {
	return "Restore file."
}

// Help implements command.
func (r *restoreFileCommand) Help() string {
	return "<RID: int> <path: string>"
}

// Execute implements command.
func (r *restoreFileCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 2 {
		return false, errors.New("expected 2 arguments")
	}

	var rid, ridError = strconv.Atoi(args.Pop())
	if ridError != nil {
		return false, ridError
	}
	var path = args.Pop()

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
	var resource, resourceError = identity.RestoreFile(ctx, (gophkeeper.ResourceID)(rid), path, vaultPassword)
	if resourceError != nil {
		return true, resourceError
	}

	fmt.Printf("Restored resource (RID: %d) to file: %s\n", rid, resource.path)

	return true, nil
}
