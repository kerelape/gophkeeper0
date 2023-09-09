package server

import (
	"context"

	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
)

// Server is gophkeeper server.
type Server struct {
	RestAddress string // the address that REST api serves at.

	Repository gophkeeper.Gophkeeper
}

var _ runnable.Runnable = (*Server)(nil)

// Run runs Server.
func (s *Server) Run(ctx context.Context) error {
	var (
		restDaemon = rest.Rest{
			Address:    s.RestAddress,
			Repository: s.Repository,
		}
	)

	var manager = runnable.NewManager()
	manager.Add(&restDaemon)
	return manager.Build().Run(ctx)
}
