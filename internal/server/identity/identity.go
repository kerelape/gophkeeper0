package identity

import "context"

// Identity is a gophkeeper's identity.
type Identity interface {
	// ComparePassword compares password agains the identity's password.
	ComparePassword(ctx context.Context, password string) error
}
