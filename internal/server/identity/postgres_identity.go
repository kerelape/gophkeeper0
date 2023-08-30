package identity

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

// PostgresIdentity is a postgres identity.
type PostgresIdentity struct {
	Connection       *pgx.Conn
	PasswordEncoding *base64.Encoding

	Username string
}

var _ Identity = (*PostgresIdentity)(nil)

// StorePiece implements Identity.
func (i *PostgresIdentity) StorePiece(ctx context.Context, piece Piece, password string) (ResourceID, error) {
	panic("unimplemented")
}

// RestorePiece implements Identity.
func (i *PostgresIdentity) RestorePiece(ctx context.Context, rid ResourceID, password string) (Piece, error) {
	panic("unimplemented")
}

// StoreBlob implements Identity.
func (i *PostgresIdentity) StoreBlob(ctx context.Context, blob Blob, password string) (ResourceID, error) {
	panic("unimplemented")
}

// RestoreBlob implements Identity.
func (i *PostgresIdentity) RestoreBlob(ctx context.Context, rid ResourceID, password string) (Blob, error) {
	panic("unimplemented")
}

// Delete implements Identity.
func (i *PostgresIdentity) Delete(context.Context, ResourceID) error {
	panic("unimplemented")
}

// List implements Identity.
func (i *PostgresIdentity) List(context.Context) ([]Resource, error) {
	panic("unimplemented")
}

func (i *PostgresIdentity) comparePassword(ctx context.Context, password string) error {
	var row = i.Connection.QueryRow(
		ctx,
		`SELECT password FROM identities WHERE username = $1`,
		i.Username,
	)
	var encodedPassword string
	if err := row.Scan(&encodedPassword); err != nil {
		var pgerr pgconn.PgError
		if errors.As(err, (any)(&pgerr)) {
			return ErrBadCredential
		}
		return err
	}

	var decodedPassword, decodePasswordError = i.PasswordEncoding.DecodeString(encodedPassword)
	if decodePasswordError != nil {
		return decodePasswordError
	}
	if err := bcrypt.CompareHashAndPassword(decodedPassword, ([]byte)(password)); err != nil {
		return errors.Join(ErrBadCredential, err)
	}
	return nil
}
