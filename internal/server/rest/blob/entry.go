package blob

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/domain"
)

// Entry is blob entry.
type Entry struct {
	Repository domain.Repository
}

// Route routes blob entry.
func (e *Entry) Route() http.Handler {
	var router = chi.NewRouter()
	router.Put("/", nil)
	router.Get("/{rid}", nil)
	return router
}
