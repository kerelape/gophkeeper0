package cli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/kerelape/gophkeeper/internal/stack"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"golang.org/x/term"
)

type registerCommand struct {
	gophkeeper gophkeeper.Gophkeeper
}

var _ command = (*registerCommand)(nil)

// Execute implements command.
func (r *registerCommand) Execute(ctx context.Context, args stack.Stack[string]) (bool, error) {
	if len(args) > 0 {
		return false, errors.New("expected 0 arguments")
	}
	var input = bufio.NewReader(os.Stdin)

	fmt.Print("Type new identity's username: ")
	var username, usernameError = input.ReadString('\n')
	if usernameError != nil {
		return true, usernameError
	}
	username = strings.TrimSuffix(username, "\n")

	fmt.Print("Type new identity's password: ")
	var password1, password1Error = term.ReadPassword((int)(syscall.Stdin))
	if password1Error != nil {
		return true, password1Error
	}

	fmt.Println()
	fmt.Print("Retype new identity's password: ")
	var password2, password2Error = term.ReadPassword((int)(syscall.Stdin))
	if password2Error != nil {
		return true, password1Error
	}
	fmt.Println()

	if !bytes.Equal(password1, password2) {
		return true, errors.New("passwords do not match")
	}

	var credential = gophkeeper.Credential{
		Username: username,
		Password: (string)(password1),
	}
	if err := r.gophkeeper.Register(ctx, credential); err != nil {
		return true, err
	}
	return true, nil
}

// Help implements command.
func (r *registerCommand) Help() string {
	return ""
}

// Description implements command.
func (*registerCommand) Description() string {
	return "Register a new identity to gophkeeper."
}
