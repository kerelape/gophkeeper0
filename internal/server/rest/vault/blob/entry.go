package blob

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is blob entry.
type Entry struct {
	Gophkeeper gophkeeper.Gophkeeper
}

// Route routes blob entry.
func (e *Entry) Route() http.Handler {
	var router = chi.NewRouter()
	router.Put("/", e.encrypt)
	router.Get("/{rid}", e.decrypt)
	return router
}

func (e *Entry) encrypt(out http.ResponseWriter, in *http.Request) {
	var token = in.Header.Get("Authorization")
	var identity, identityError = e.Gophkeeper.Identity(in.Context(), (gophkeeper.Token)(token))
	if identityError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(identityError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	var password = in.Header.Get("X-Password")
	if password == "" {
		var status = http.StatusUnauthorized
		http.Error(out, http.StatusText(status), status)
		return
	}

	var blob = gophkeeper.Blob{
		Meta:    in.Header.Get("X-Meta"),
		Content: in.Body,
	}
	rid, storeError := identity.StoreBlob(in.Context(), blob, password)
	if storeError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(storeError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	out.WriteHeader(http.StatusCreated)
	var response struct {
		RID int64 `json:"rid"`
	}
	response.RID = (int64)(rid)
	if err := json.NewEncoder(out).Encode(response); err != nil {
		log.Printf("Failed to write response: %s", err.Error())
	}
}

func (e *Entry) decrypt(out http.ResponseWriter, in *http.Request) {
	var token = in.Header.Get("Authorization")
	var identity, identityError = e.Gophkeeper.Identity(in.Context(), (gophkeeper.Token)(token))
	if identityError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(identityError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	var rid, ridError = strconv.Atoi(chi.URLParam(in, "rid"))
	if ridError != nil {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	var password = in.Header.Get("X-Password")
	if password == "" {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	var blob, restoreError = identity.RestoreBlob(in.Context(), (gophkeeper.ResourceID)(rid), password)
	if restoreError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(identityError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}
	defer blob.Content.Close()

	out.Header().Set("Content-Type", "application/octet-stream")
	out.Header().Set("Content-Disposition", "attachment")
	out.Header().Set("X-Meta", blob.Meta)
	out.WriteHeader(http.StatusOK)

	var output = bufio.NewWriter(out)
	if _, err := output.ReadFrom(blob.Content); err != nil {
		log.Printf("failed to write content: %s", err.Error())
	}
	if err := output.Flush(); err != nil {
		log.Printf("failed to flush content: %s", err.Error())
	}
}
