package piece

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/domain"
)

// Entry is piece entry.
type Entry struct {
	Repository domain.Repository
}

// Route routes piece entry.
func (e *Entry) Route() http.Handler {
	var router = chi.NewRouter()
	router.Put("/", e.encrypt)
	router.Get("/{rid}", e.decrypt)
	return router
}

func (e *Entry) encrypt(out http.ResponseWriter, in *http.Request) {

}

func (e *Entry) decrypt(out http.ResponseWriter, in *http.Request) {

}
