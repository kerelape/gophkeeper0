package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type restoreCardCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*restoreCardCommand)(nil)

// Description implements command.
func (r *restoreCardCommand) Description() string {
	return "Restore card information."
}

// Help implements command.
func (r *restoreCardCommand) Help() string {
	return "<RID: int>"
}

// Execute implements command.
func (r *restoreCardCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
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
	var resource, resourceError = identity.RestoreCard(ctx, (gophkeeper.ResourceID)(rid), vaultPassword)
	if resourceError != nil {
		return true, resourceError
	}

	fmt.Printf("(%d) Card\n", rid)
	fmt.Printf("\nCard Number\n%s\n", resource.ccn)
	fmt.Printf("\nExpiration Date    CVV\n")
	fmt.Printf("%s              %s\n", resource.exp, resource.cvv)
	fmt.Printf("\nCard Holder\n%s\n", resource.holder)

	return true, nil
}
