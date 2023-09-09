package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
	"golang.org/x/crypto/acme/autocert"
)

// Rest is gophkeeper's rest API.
type Rest struct {
	Address string // Address is the address that REST api serves at.

	Repository gophkeeper.Gophkeeper
}

var _ runnable.Runnable = (*Rest)(nil)

// Run implement runnable.Runnable for Rest.
func (r *Rest) Run(ctx context.Context) error {
	var (
		entry = Entry{
			Gophkeeper: r.Repository,
		}
		server = http.Server{
			Addr:    r.Address,
			Handler: entry.Route(),
		}
		serverErrorChannel = make(chan error)
	)
	go func() {
		serverErrorChannel <- server.Serve(autocert.NewListener())
	}()
	select {
	case <-ctx.Done():
		return server.Shutdown(ctx)
	case err := <-serverErrorChannel:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
