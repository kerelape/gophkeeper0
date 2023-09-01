package domain

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"os"
	"path"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
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
	BlobsDir         string

	Username string
}

var _ gophkeeper.Identity = (*PostgresIdentity)(nil)

// StorePiece implements Identity.
func (i *PostgresIdentity) StorePiece(ctx context.Context, piece gophkeeper.Piece, password string) (gophkeeper.ResourceID, error) {
	var (
		salt []byte = make([]byte, 8)
		iv   []byte = make([]byte, keyLen)
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
	return (gophkeeper.ResourceID)(rid), nil
}

// RestorePiece implements Identity.
func (i *PostgresIdentity) RestorePiece(ctx context.Context, rid gophkeeper.ResourceID, password string) (gophkeeper.Piece, error) {
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
			return gophkeeper.Piece{}, gophkeeper.ErrResourceNotFound
		}
		return gophkeeper.Piece{}, err
	}
	var key = pbkdf2.Key(([]byte)(password), salt, keyIter, keyLen, sha256.New)
	var block, blockError = aes.NewCipher(key)
	if blockError != nil {
		return gophkeeper.Piece{}, blockError
	}
	var aesgcm, aesgcmError = cipher.NewGCM(block)
	if aesgcmError != nil {
		return gophkeeper.Piece{}, aesgcmError
	}
	var decryptedContent, openError = aesgcm.Open(nil, iv, content, nil)
	if openError != nil {
		return gophkeeper.Piece{}, openError
	}
	var piece = gophkeeper.Piece{
		Meta:    meta,
		Content: decryptedContent,
	}
	return piece, nil
}

// StoreBlob implements Identity.
func (i *PostgresIdentity) StoreBlob(ctx context.Context, blob gophkeeper.Blob, password string) (gophkeeper.ResourceID, error) {
	defer blob.Content.Close()
	var (
		salt     []byte = make([]byte, 8)
		iv       []byte = make([]byte, keyLen)
		key      []byte
		rid      int64
		blobPath string
	)
	if _, err := rand.Read(salt); err != nil {
		return -1, err
	}
	if _, err := rand.Read(iv); err != nil {
		return -1, err
	}
	key = pbkdf2.Key(([]byte)(password), salt, keyIter, keyLen, sha256.New)

	var transaction, beginError = i.Connection.Begin(ctx)
	if beginError != nil {
		return -1, beginError
	}
	var ridRow = transaction.QueryRow(
		ctx,
		`INSERT INTO blobs(owner, meta, salt, iv) VALUES($1, $2, $3, $4) RETURNING rid`,
		i.Username,
		blob.Meta,
		salt, iv,
	)
	if err := ridRow.Scan(&rid); err != nil {
		if err := transaction.Rollback(ctx); err != nil {
			return -1, err
		}
		return -1, err
	}
	blobPath = path.Join(i.BlobsDir, strconv.Itoa((int)(rid)))
	_, updateError := transaction.Exec(
		ctx,
		`UPDATE blobs SET path = $1 WHERE rid = $2`,
		blobPath, rid,
	)
	if updateError != nil {
		return -1, updateError
	}
	if err := transaction.Commit(ctx); err != nil {
		return -1, err
	}

	var file, createError = os.Create(blobPath)
	if createError != nil {
		return -1, createError
	}
	var block, blockError = aes.NewCipher(key)
	if blockError != nil {
		return -1, blockError
	}
	var writer = cipher.StreamWriter{
		S: cipher.NewCTR(block, iv),
		W: file,
	}
	var reader = bufio.NewReader(blob.Content)
	if _, err := reader.WriteTo(writer); err != nil {
		return -1, err
	}
	return (gophkeeper.ResourceID)(rid), nil
}

// RestoreBlob implements Identity.
func (i *PostgresIdentity) RestoreBlob(ctx context.Context, rid gophkeeper.ResourceID, password string) (gophkeeper.Blob, error) {
	var (
		owner    string
		meta     string
		blobPath string
		iv       []byte
		salt     []byte
		key      []byte
	)
	var row = i.Connection.QueryRow(
		ctx,
		`SELECT (owner, meta, path, iv, salt) FROM blobs WHERE rid = $1`,
		(int64)(rid),
	)
	if err := row.Scan(&owner, &meta, &blobPath, &iv, &salt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return gophkeeper.Blob{}, gophkeeper.ErrResourceNotFound
		}
		return gophkeeper.Blob{}, err
	}
	if owner != i.Username {
		return gophkeeper.Blob{}, gophkeeper.ErrResourceNotFound
	}

	var file, openError = os.Open(blobPath)
	if openError != nil {
		return gophkeeper.Blob{}, openError
	}
	key = pbkdf2.Key(([]byte)(password), salt, keyIter, keyLen, sha256.New)
	var block, blockError = aes.NewCipher(key)
	if blockError != nil {
		return gophkeeper.Blob{}, blockError
	}
	var blob = gophkeeper.Blob{
		Meta: meta,
		Content: &readCloser{
			reader: cipher.StreamReader{
				S: cipher.NewCTR(block, iv),
				R: file,
			},
			closer: file,
		},
	}
	return blob, nil
}

// Delete implements Identity.
func (i *PostgresIdentity) Delete(context.Context, gophkeeper.ResourceID) error {
	panic("unimplemented")
}

// List implements Identity.
func (i *PostgresIdentity) List(context.Context) ([]gophkeeper.Resource, error) {
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
			return gophkeeper.ErrBadCredential
		}
		return err
	}

	var decodedPassword, decodePasswordError = i.PasswordEncoding.DecodeString(encodedPassword)
	if decodePasswordError != nil {
		return decodePasswordError
	}
	if err := bcrypt.CompareHashAndPassword(decodedPassword, ([]byte)(password)); err != nil {
		return errors.Join(gophkeeper.ErrBadCredential, err)
	}
	return nil
}

type readCloser struct {
	reader io.Reader
	closer io.Closer
}

func (rc *readCloser) Read(p []byte) (int, error) {
	return rc.reader.Read(p)
}

func (rc *readCloser) Close() error {
	return rc.closer.Close()
}
