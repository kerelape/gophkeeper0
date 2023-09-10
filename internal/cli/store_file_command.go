package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type storeFileCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*storeFileCommand)(nil)

// Description implements command.
func (s *storeFileCommand) Description() string {
	return "Stores a file."
}

// Help implements command.
func (s *storeFileCommand) Help() string {
	return "<path: string>"
}

// Execute implements command.
func (s *storeFileCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) != 1 {
		return false, errors.New("expected 1 argument")
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

	var (
		identity = identity{
			origin: gophkeeperIdentity,
		}
		resource = fileResource{
			description: description,
			path:        args.Pop(),
		}
	)
	var rid, ridError = identity.StoreFile(ctx, resource, vaultPassword)
	if ridError != nil {
		return true, ridError
	}

	fmt.Printf("Successfully stored file.\n")
	fmt.Printf("RID of the newly stored resource is %d.\n", rid)

	return true, nil
}
