package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type storeTextCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*storeTextCommand)(nil)

// Description implements command.
func (s *storeTextCommand) Description() string {
	return "Store a text note."
}

// Help implements command.
func (s *storeTextCommand) Help() string {
	return ""
}

// Execute implements command.
func (s *storeTextCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) > 0 {
		return false, errors.New("expected 0 arguments")
	}

	var gophkeeperIdentity, gophkeeperIdentityError = authenticate(ctx, s.gophkeeper)
	if gophkeeperIdentityError != nil {
		return true, gophkeeperIdentityError
	}

	var vaultPassword, vaultPasswordError = vaultPassword(ctx)
	if vaultPasswordError != nil {
		return true, vaultPasswordError
	}

	var description, descriptionError = description(ctx)
	if descriptionError != nil {
		return true, descriptionError
	}

	var content, contentError = text(ctx)
	if contentError != nil {
		return true, contentError
	}

	var (
		identity = identity{
			origin: gophkeeperIdentity,
		}
		resource = textResource{
			description: description,
			content:     content,
		}
	)

	var rid, ridError = identity.StoreText(ctx, resource, vaultPassword)
	if ridError != nil {
		return true, ridError
	}

	fmt.Printf("\nSuccessfully saved text note.\n")
	fmt.Printf("RID of the newly stored resource is %d.\n", rid)

	return true, nil
}
