package server

import (
	"context"

	"github.com/kerelape/gophkeeper/internal/server/domain"
	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

// Server is gophkeeper server.
type Server struct {
	RestAddress  string // the address that REST api serves at.
	RestCertFile string // path to the cert file used for REST api.
	RestKeyFile  string // path to the key file used for REST api.

	Repository domain.Repository
}

var _ runnable.Runnable = (*Server)(nil)

// Run runs Server.
func (s *Server) Run(ctx context.Context) error {
	var (
		restDaemon = rest.Rest{
			Address:    s.RestAddress,
			CertFile:   s.RestCertFile,
			KeyFile:    s.RestKeyFile,
			Repository: s.Repository,
		}
	)

	var manager = runnable.NewManager()
	manager.Add(&restDaemon)
	return manager.Build().Run(ctx)
}
