package rest

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
	"github.com/pior/runnable"
	"golang.org/x/crypto/acme/autocert"
)

// Rest is gophkeeper's rest API.
type Rest struct {
	Address       string // Address is the address that REST api serves at.
	UseTLS        bool
	HostWhilelist []string

	Gophkeeper gophkeeper.Gophkeeper
}

var _ runnable.Runnable = (*Rest)(nil)

// Run implement runnable.Runnable for Rest.
func (r *Rest) Run(ctx context.Context) error {
	var (
		entry = Entry{
			Gophkeeper: r.Gophkeeper,
		}
		server = http.Server{
			Addr:    r.Address,
			Handler: entry.Route(),
		}
		serverErrorChannel = make(chan error)
	)
	go func() {
		if r.UseTLS {
			var certmanager = autocert.Manager{
				Cache:      autocert.DirCache("cache"),
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(r.HostWhilelist...),
			}
			server.TLSConfig = certmanager.TLSConfig()
			serverErrorChannel <- server.ListenAndServeTLS("", "")
		} else {
			log.Printf("\033[31m!CONNECTION IS NOT SECURED!\033[0m")
			serverErrorChannel <- server.ListenAndServe()
		}
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
