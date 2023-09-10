package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type restoreCredentialCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*restoreCredentialCommand)(nil)

// Description implements command.
func (r *restoreCredentialCommand) Description() string {
	return "Restore a username-password pair."
}

// Help implements command.
func (r *restoreCredentialCommand) Help() string {
	return "<RID: int>"
}

// Execute implements command.
func (r *restoreCredentialCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("expected 1 arguments")
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
	var resource, resourceError = identity.RestoreCredential(ctx, (gophkeeper.ResourceID)(rid), vaultPassword)
	if resourceError != nil {
		return true, resourceError
	}

	fmt.Printf("(%d) Credential\n", rid)
	fmt.Printf("\tUsername: %s\n", resource.username)
	fmt.Printf("\tPassword: %s\n\n", resource.password)

	return true, nil
}
