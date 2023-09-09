package cli

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type listCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*listCommand)(nil)

// Description implements command.
func (l *listCommand) Description() string {
	return "List out all available resources."
}

// Help implements command.
func (l *listCommand) Help() string {
	return ""
}

// Execute implements command.
func (l *listCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) > 0 {
		return false, errors.New("expected 0 arguments")
	}
	var gophkeeperIdentity, identityError = authenticate(ctx, l.gophkeeper)
	if identityError != nil {
		return true, identityError
	}
	var identity = identity{
		origin: gophkeeperIdentity,
	}
	var resources, resourcesError = identity.List(ctx)
	if resourcesError != nil {
		return true, resourcesError
	}
	fmt.Printf("%d resources found\n", len(resources))
	for _, r := range resources {
		fmt.Printf(
			"(RID: %d)\n\tType: %s\n\tDescription: %s\n",
			r.RID,
			r.Type.String(),
			strings.ReplaceAll(r.Description, "\n", " "),
		)
	}
	return true, nil
}
