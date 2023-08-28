package server

import (
	"context"

	"github.com/kerelape/gophkeeper/internal/server/rest"
	"github.com/pior/runnable"
)

// Server is gophkeeper server.
type Server struct {
	Rest rest.Rest
}

var _ runnable.Runnable = (*Server)(nil)

// Run runs Server.
func (s *Server) Run(ctx context.Context) error {
	var manager = runnable.NewManager()
	manager.Add(&s.Rest)
	return manager.Build().Run(ctx)
}
