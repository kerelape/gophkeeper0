package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kerelape/gophkeeper/internal/deferred"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
	"golang.org/x/crypto/bcrypt"
)

// go:embed initialize_postgres.sql
var initializePostgresQuery string

// PostgresRepository is a postgresql identity repository.
type PostgresRepository struct {
	connection deferred.Deferred[*pgx.Conn]

	PasswordEncoding *base64.Encoding

	DSN           string
	TokenSecret   []byte
	TokenLifespan time.Duration
	BlobsDir      string

	UsernameMinLength uint
	PasswordMinLength uint
}

var (
	_ gophkeeper.Gophkeeper = (*PostgresRepository)(nil)
	_ runnable.Runnable     = (*PostgresRepository)(nil)
)

// Register implements Repository.
func (r *PostgresRepository) Register(ctx context.Context, credential gophkeeper.Credential) error {
	var connection, connectionError = r.connection.Get(ctx)
	if connectionError != nil {
		return connectionError
	}

	if len(credential.Username) < (int)(r.UsernameMinLength) {
		return gophkeeper.ErrBadCredential
	}
	if len(credential.Password) < (int)(r.PasswordMinLength) {
		return gophkeeper.ErrBadCredential
	}

	var password, passwordError = bcrypt.GenerateFromPassword(
		([]byte)(credential.Password),
		bcrypt.DefaultCost,
	)
	if passwordError != nil {
		return passwordError
	}

	_, insertError := connection.Exec(
		ctx,
		`INSERT INTO identities(username, password) VALUES($1, $2)`,
		credential.Username,
		r.PasswordEncoding.EncodeToString(password),
	)
	if insertError != nil {
		if err := new(pgconn.PgError); errors.As(insertError, &err) && err.Code == "23505" {
			return gophkeeper.ErrIdentityDuplicate
		}
		return insertError
	}

	return nil
}

// Authenticate implements Repository.
func (r *PostgresRepository) Authenticate(ctx context.Context, credential gophkeeper.Credential) (gophkeeper.Token, error) {
	var connection, connectionError = r.connection.Get(ctx)
	if connectionError != nil {
		return (gophkeeper.Token)(""), connectionError
	}

	var identity = PostgresIdentity{
		Connection:       connection,
		PasswordEncoding: r.PasswordEncoding,
		Username:         credential.Username,
	}
	if err := identity.comparePassword(ctx, credential.Password); err != nil {
		return (gophkeeper.Token)(""), err
	}

	var rawToken = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": time.Now().Add(r.TokenLifespan).Unix(),
			"sub": credential.Username,
		},
	)
	var token, signTokenError = rawToken.SignedString(r.TokenSecret)
	if signTokenError != nil {
		return (gophkeeper.Token)(""), signTokenError
	}
	return (gophkeeper.Token)(token), nil
}

// Identity implements Repository.
func (r *PostgresRepository) Identity(ctx context.Context, token gophkeeper.Token) (gophkeeper.Identity, error) {
	var parsedToken, parseTokenError = jwt.Parse(
		(string)(token),
		func(t *jwt.Token) (interface{}, error) {
			return r.TokenSecret, nil
		},
	)
	if parseTokenError != nil {
		return nil, gophkeeper.ErrBadCredential
	}

	var claims = parsedToken.Claims
	if exp, err := claims.GetExpirationTime(); err == nil {
		if exp.Before(time.Now()) {
			return nil, gophkeeper.ErrBadCredential
		}
	} else {
		return nil, gophkeeper.ErrBadCredential
	}

	var username string
	if sub, err := claims.GetSubject(); err == nil {
		username = sub
	} else {
		return nil, gophkeeper.ErrBadCredential
	}

	var connection, connectionError = r.connection.Get(ctx)
	if connectionError != nil {
		return nil, connectionError
	}

	var identity = &PostgresIdentity{
		Connection:       connection,
		PasswordEncoding: r.PasswordEncoding,
		Username:         username,
		BlobsDir:         r.BlobsDir,
	}
	return identity, nil
}

// Run implements Runnable.
func (r *PostgresRepository) Run(ctx context.Context) error {
	var connection, connectError = pgx.Connect(ctx, r.DSN)
	if connectError != nil {
		return connectError
	}
	defer connection.Close(context.Background())

	_, initializeError := connection.Exec(ctx, initializePostgresQuery)
	if initializeError != nil {
		return initializeError
	}

	r.connection.Set(connection)

	<-ctx.Done()
	return ctx.Err()
}
