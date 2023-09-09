package cli

import (
	"context"

	"github.com/kerelape/gophkeeper/internal/stack"
)

type command interface {
	Help() string
	Description() string
	Execute(ctx context.Context, args stack.Stack[string]) (bool, error)
}
