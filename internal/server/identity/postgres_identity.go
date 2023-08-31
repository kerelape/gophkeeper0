package identity

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLen  = 32
	keyIter = 4096
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
	var (
		salt []byte = make([]byte, 8)
		iv   []byte = make([]byte, 12)
		key  []byte
	)
	if _, err := rand.Read(salt); err != nil {
		return -1, err
	}
	if _, err := rand.Read(iv); err != nil {
		return -1, err
	}
	key = pbkdf2.Key(([]byte)(password), salt, keyIter, keyLen, sha256.New)
	var block, blockError = aes.NewCipher(key)
	if blockError != nil {
		return -1, blockError
	}
	var aesgcm, aesgcmError = cipher.NewGCM(block)
	if aesgcmError != nil {
		return -1, aesgcmError
	}
	var content = aesgcm.Seal(nil, iv, piece.Content, nil)

	result := i.Connection.QueryRow(
		ctx,
		`INSERT INTO pieces(owner, meta, content, salt, iv) VALUES($1, $2, $3, $4, $5) RETURNING rid`,
		i.Username,
		piece.Meta, content,
		salt, iv,
	)
	var rid int64
	if err := result.Scan(&rid); err != nil {
		return -1, err
	}
	return (ResourceID)(rid), nil
}

// RestorePiece implements Identity.
func (i *PostgresIdentity) RestorePiece(ctx context.Context, rid ResourceID, password string) (Piece, error) {
	var (
		owner   string
		meta    string
		content []byte
		iv      []byte
		salt    []byte
	)
	var result = i.Connection.QueryRow(
		ctx,
		`SELECT (owner, meta, content, iv, salt) FROM pieces WHERE rid = $1`,
		(int64)(rid),
	)
	if err := result.Scan(&owner, &meta, &content, &iv, &salt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Piece{}, ErrResourceNotFound
		}
		return Piece{}, err
	}
	var key = pbkdf2.Key(([]byte)(password), salt, keyIter, keyLen, sha256.New)
	var block, blockError = aes.NewCipher(key)
	if blockError != nil {
		return Piece{}, blockError
	}
	var aesgcm, aesgcmError = cipher.NewGCM(block)
	if aesgcmError != nil {
		return Piece{}, aesgcmError
	}
	var decryptedContent, openError = aesgcm.Open(nil, iv, content, nil)
	if openError != nil {
		return Piece{}, openError
	}
	var piece = Piece{
		Meta:    meta,
		Content: decryptedContent,
	}
	return piece, nil
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
