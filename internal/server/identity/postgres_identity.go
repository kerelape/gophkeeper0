package identity

import (
	"github.com/jackc/pgx/v5"
)

// PostgresIdentity is a postgres identity.
type PostgresIdentity struct {
	Connection *pgx.Conn
	Username   string
}

var _ Identity = (*PostgresIdentity)(nil)
