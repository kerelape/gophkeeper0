package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/kerelape/gophkeeper/internal/server/identity"
	"github.com/pior/runnable"
)

// Rest is gophkeeper's rest API.
type Rest struct {
	Address  string // Address is the address that REST api serves at.
	CertFile string // CertFile is path to the cert file used for REST api.
	KeyFile  string // KeyFile is path to the key file used for REST api.

	Repository identity.Repository
}

var _ runnable.Runnable = (*Rest)(nil)

// Run implement runnable.Runnable for Rest.
func (r *Rest) Run(ctx context.Context) error {
	var (
		entry = Entry{
			Repository: r.Repository,
		}
		server = http.Server{
			Addr:    r.Address,
			Handler: entry.Route(),
		}
		serverErrorChannel = make(chan error)
	)
	go func() {
		serverErrorChannel <- server.ListenAndServeTLS(r.CertFile, r.KeyFile)
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
