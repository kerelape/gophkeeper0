package vault

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/internal/server/rest/vault/blob"
	"github.com/kerelape/gophkeeper/internal/server/rest/vault/piece"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is vault entry.
type Entry struct {
	Gophkeeper gophkeeper.Gophkeeper
}

// Route routes vault entry.
func (e *Entry) Route() http.Handler {
	var (
		piece = piece.Entry{
			Gophkeeper: e.Gophkeeper,
		}
		blob = blob.Entry{
			Gophkeeper: e.Gophkeeper,
		}
	)
	var router = chi.NewRouter()
	router.Mount("/piece", piece.Route())
	router.Mount("/blob", blob.Route())
	router.Get("/", nil)
	router.Delete("/{rid}", nil)
	return router
}

// @todo #32 Implement `GET /vault`
func (e *Entry) get(out http.ResponseWriter, in *http.Request) {
	panic("unimplemented")
}

// @todo #32 Implement `DELETE /vault/{rid}`
func (e *Entry) delete(out http.ResponseWriter, in *http.Request) {
	panic("unimplemented")
}
