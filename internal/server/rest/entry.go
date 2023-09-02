package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/rest/login"
	"github.com/kerelape/gophkeeper/internal/server/rest/register"
	"github.com/kerelape/gophkeeper/internal/server/rest/vault"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is the REST api entry.
type Entry struct {
	Gophkeeper gophkeeper.Gophkeeper
}

// Route routes Entry into an http.Handler.
func (e *Entry) Route() http.Handler {
	var (
		register = register.Entry{
			Gophkeeper: e.Gophkeeper,
		}
		login = login.Entry{
			Gophkeeper: e.Gophkeeper,
		}
		vault = vault.Entry{
			Gophkeeper: e.Gophkeeper,
		}
	)
	var router = chi.NewRouter()
	router.Mount("/register", register.Route())
	router.Mount("/login", login.Route())
	router.Mount("/vault", vault.Route())
	return router
}
