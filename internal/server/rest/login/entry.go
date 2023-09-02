package login

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is login entry.
type Entry struct {
	Gophkeeper gophkeeper.Gophkeeper
}

// Route routes this entry into an http.Handler.
func (e *Entry) Route() http.Handler {
	var router = chi.NewRouter()
	router.Post("/", e.post)
	return router
}

func (e *Entry) post(out http.ResponseWriter, in *http.Request) {
	var requestBody map[string]any
	if err := json.NewDecoder(in.Body).Decode(&requestBody); err != nil {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	var credential gophkeeper.Credential
	if val, ok := requestBody["username"].(string); ok {
		credential.Username = val
	} else {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	if val, ok := requestBody["password"].(string); ok {
		credential.Password = val
	} else {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	var token, authenticateError = e.Gophkeeper.Authenticate(in.Context(), credential)
	if authenticateError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(authenticateError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	out.Header().Set("Authorization", (string)(token))
	out.WriteHeader(http.StatusOK)
}
