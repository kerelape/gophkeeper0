package server

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/kerelape/gophkeeper/internal/server/postgres"
	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

// Server is gophkeeper server.
type Server struct {
	RestAddress       string // the address that REST api serves at.
	RestUseTLS        bool
	RestHostWhilelist []string

	DatabaseDSN   string
	BlobsDir      string
	TokenSecret   []byte
	TokenLifespan time.Duration

	UsernameMinLength uint
	PasswordMinLength uint
}

var _ runnable.Runnable = (*Server)(nil)

// Run runs Server.
func (s *Server) Run(ctx context.Context) error {
	var (
		gophkeeper = postgres.Gophkeeper{
			PasswordEncoding: base64.RawStdEncoding,

			DSN:      s.DatabaseDSN,
			BlobsDir: s.BlobsDir,

			TokenSecret:   s.TokenSecret,
			TokenLifespan: s.TokenLifespan,

			UsernameMinLength: s.UsernameMinLength,
			PasswordMinLength: s.PasswordMinLength,
		}
		restDaemon = rest.Rest{
			Address:       s.RestAddress,
			Gophkeeper:    &gophkeeper,
			UseTLS:        s.RestUseTLS,
			HostWhilelist: s.RestHostWhilelist,
		}
	)

	var manager = runnable.NewManager()
	manager.Add(&gophkeeper)
	manager.Add(&restDaemon)
	return manager.Build().Run(ctx)
}
