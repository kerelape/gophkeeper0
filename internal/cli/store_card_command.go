package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type storeCardCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*storeCardCommand)(nil)

// Description implements command.
func (s *storeCardCommand) Description() string {
	return "Store card information."
}

// Help implements command.
func (s *storeCardCommand) Help() string {
	return ""
}

// Execute implements command.
func (s *storeCardCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 0 {
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

	var card, cardError = cardCredential(ctx)
	if cardError != nil {
		return true, cardError
	}

	var (
		identity = identity{
			origin: gophkeeperIdentity,
		}
		resource = cardResource{
			cardInfo:    card,
			description: description,
		}
	)
	var rid, ridError = identity.StoreCard(ctx, resource, vaultPassword)
	if ridError != nil {
		return true, ridError
	}

	fmt.Printf("Successfully stored card.\n")
	fmt.Printf("RID of the newly stored resource is %d.\n", rid)

	return true, nil
}
