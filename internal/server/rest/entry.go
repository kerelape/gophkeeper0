package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/rest/blob"
	"github.com/kerelape/gophkeeper/internal/server/rest/login"
	"github.com/kerelape/gophkeeper/internal/server/rest/piece"
	"github.com/kerelape/gophkeeper/internal/server/rest/register"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is the REST api entry.
type Entry struct {
	Repository gophkeeper.Repository
}

// Route routes Entry into an http.Handler.
func (e *Entry) Route() http.Handler {
	var (
		registretion = register.Entry{
			Repository: e.Repository,
		}
		authentication = login.Entry{
			Repository: e.Repository,
		}
		piece = piece.Entry{
			Repository: e.Repository,
		}
		blob = blob.Entry{
			Repository: e.Repository,
		}
	)
	var router = chi.NewRouter()
	router.Mount("/register", registretion.Route())
	router.Mount("/login", authentication.Route())
	router.Mount("/piece", piece.Route())
	router.Mount("/blob", blob.Route())
	return router
}
