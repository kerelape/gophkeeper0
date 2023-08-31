package domain

import (
	"context"
	"errors"
)

var (
	// ErrBadCredential indecates that the credential provided are bad.
	ErrBadCredential = errors.New("bad credential")

	// ErrInvalidToken indecates that the token provided is invalid.
	ErrInvalidToken = errors.New("invalid token")

	// ErrIdentityDuplicate indecates that there is already such an identity.
	ErrIdentityDuplicate = errors.New("identity already exists")
)

type (
	// Token is a JWT token.
	Token string

	// Credential is a user authentication credential.
	Credential struct {
		Username string
		Password string
	}
)

// Repository is an identity reposity.
type Repository interface {
	// Register registers a new identity into the reposity.
	Register(context.Context, Credential) error

	// Authenticate authenticates an identity and returns an access token.
	Authenticate(context.Context, Credential) (Token, error)

	// Identity returns the identity associated with the token.
	Identity(context.Context, Token) (Identity, error)
}
