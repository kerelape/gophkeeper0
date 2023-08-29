package identity

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/kerelape/gophkeeper/internal/deferred"
	"github.com/pior/runnable"
	"golang.org/x/crypto/bcrypt"
)

// go:embed initialize_postgres.sql
var initializePostgresQuery string

// PostgresRepository is a postgresql identity repository.
type PostgresRepository struct {
	connection deferred.Deferred[*pgx.Conn]

	DSN       string
	JWTSecret string

	UsernameMinLength uint
	PasswordMinLength uint
}

var (
	_ Repository        = (*PostgresRepository)(nil)
	_ runnable.Runnable = (*PostgresRepository)(nil)
)

// Register implements Repository.
func (r *PostgresRepository) Register(ctx context.Context, credential Credential) error {
	var connection, connectionError = r.connection.Get(ctx)
	if connectionError != nil {
		return connectionError
	}

	if len(credential.Username) < (int)(r.UsernameMinLength) {
		return ErrBadCredential
	}
	if len(credential.Password) < (int)(r.PasswordMinLength) {
		return ErrBadCredential
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
		base64.StdEncoding.EncodeToString(password),
	)
	if insertError != nil {
		if err := new(pgconn.PgError); errors.As(insertError, &err) && err.Code == "23505" {
			return ErrIdentityDuplicate
		}
		return insertError
	}

	return nil
}

// Authenticate implements Repository.
func (r *PostgresRepository) Authenticate(ctx context.Context, credential Credential) (Token, error) {
	panic("unimplemented")
}

// Identity implements Repository.
func (r *PostgresRepository) Identity(ctx context.Context, token Token) (Identity, error) {
	panic("unimplemented")
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
