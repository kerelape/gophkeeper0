package server

import "context"

// Server is gophkeeper server.
//
// @todo #2 Implement HTTPS interface.
type Server struct{}

// MakeServer returns a new Server.
func MakeServer() Server {
	return Server{}
}

// Run runs Server.
func (s *Server) Run(ctx context.Context) error {
	return nil
}
