package gophkeeper

import (
	"context"
	"net/http"
)

// RestGophkeeper is a remote gophkeeper.
type RestGophkeeper struct {
	Client http.Client
	Server string
}

var _ Gophkeeper = (*RestGophkeeper)(nil)

// Register implements Gophkeeper.
//
// @todo #28 Implement registration.
func (*RestGophkeeper) Register(context.Context, Credential) error {
	panic("unimplemented")
}

// Authenticate implements Gophkeeper.
//
// @todo #28 Implement authentication.
func (*RestGophkeeper) Authenticate(context.Context, Credential) (Token, error) {
	panic("unimplemented")
}

// Identity implements Gophkeeper.
//
// @todo #28 Implement Identity.
func (*RestGophkeeper) Identity(context.Context, Token) (Identity, error) {
	panic("unimplemented")
}
