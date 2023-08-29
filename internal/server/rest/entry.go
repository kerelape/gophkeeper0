package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/identity"
	"github.com/kerelape/gophkeeper/internal/server/rest/register"
)

// Entry is the REST api entry.
//
// @todo #7 Implement authentication.
// @todo #7 Implement text storage.
// @todo #7 Implement blob storage.
// @todo #7 Implement bank card storage.
type Entry struct {
	Repository identity.Repository
}

// Route routes Entry into an http.Handler.
func (e *Entry) Route() http.Handler {
	var router = chi.NewRouter()
	var registretion = register.Entry{
		Repository: e.Repository,
	}
	router.Mount("register", registretion.Route())
	return router
}
