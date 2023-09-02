package postgres

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"
	"path"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	composedreadcloser "github.com/kerelape/gophkeeper/internal/composed_read_closer"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLen  = 32
	keyIter = 4096
)

// Identity is a postgres identity.
type Identity struct {
	Connection       *pgx.Conn
	PasswordEncoding *base64.Encoding
	BlobsDir         string

	Username string
}

var _ gophkeeper.Identity = (*Identity)(nil)

// StorePiece implements Identity.
func (i *Identity) StorePiece(ctx context.Context, piece gophkeeper.Piece, password string) (gophkeeper.ResourceID, error) {
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

	var transaction, transactionError = i.Connection.Begin(ctx)
	if transactionError != nil {
		return -1, transactionError
	}
	defer transaction.Rollback(context.Background())

	insertPieceResult := transaction.QueryRow(
		ctx,
		`INSERT INTO pieces(content, salt, iv) VALUES($1, $2, $3) RETURNING id`,
		content, salt, iv,
	)
	var id int
	if err := insertPieceResult.Scan(&id); err != nil {
		return -1, err
	}
	insertResourceResult := transaction.QueryRow(
		ctx,
		`INSERT INTO resources(meta, resource, type) VALUES($1, $2, $3) RETURNING id`,
		piece.Meta, id, (int)(gophkeeper.ResourceTypePiece),
	)
	var rid int64
	if err := insertResourceResult.Scan(&rid); err != nil {
		return -1, err
	}
	if err := transaction.Commit(ctx); err != nil {
		return -1, err
	}

	return (gophkeeper.ResourceID)(rid), nil
}

// RestorePiece implements Identity.
func (i *Identity) RestorePiece(ctx context.Context, rid gophkeeper.ResourceID, password string) (gophkeeper.Piece, error) {
	var (
		meta    string
		content []byte
		iv      []byte
		salt    []byte
	)

	var queryResourceResult = i.Connection.QueryRow(
		ctx,
		`SELECT (meta, resource) FROM resources WHERE id = $1 AND owner = $2 AND type = $3`,
		(int64)(rid), i.Username, (int)(gophkeeper.ResourceTypePiece),
	)
	var id int
	if err := queryResourceResult.Scan(&meta, &id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return gophkeeper.Piece{}, gophkeeper.ErrResourceNotFound
		}
		return gophkeeper.Piece{}, err
	}
	var queryPieceResult = i.Connection.QueryRow(
		ctx,
		`SELECT (content, iv, salt) FROM pieces WHERE id = $1`,
		id,
	)
	if err := queryPieceResult.Scan(&content, &iv, &salt); err != nil {
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
func (i *Identity) StoreBlob(ctx context.Context, blob gophkeeper.Blob, password string) (gophkeeper.ResourceID, error) {
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
func (i *Identity) RestoreBlob(ctx context.Context, rid gophkeeper.ResourceID, password string) (gophkeeper.Blob, error) {
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
		Content: &composedreadcloser.ComposedReadCloser{
			Reader: cipher.StreamReader{
				S: cipher.NewCTR(block, iv),
				R: file,
			},
			Closer: file,
		},
	}
	return blob, nil
}

// Delete implements Identity.
func (i *Identity) Delete(context.Context, gophkeeper.ResourceID) error {
	panic("unimplemented")
}

// List implements Identity.
func (i *Identity) List(context.Context) ([]gophkeeper.Resource, error) {
	panic("unimplemented")
}

func (i *Identity) comparePassword(ctx context.Context, password string) error {
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
