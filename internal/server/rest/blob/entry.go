package blob

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

// Entry is blob entry.
type Entry struct {
	Repository gophkeeper.Repository
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
	var identity, identityError = e.Repository.Identity(in.Context(), (gophkeeper.Token)(token))
	if identityError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(identityError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	var body = bufio.NewReader(in.Body)
	var password, passwordError = body.ReadString(0x00)
	if passwordError != nil {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}
	password, _ = strings.CutSuffix(password, (string)(0x00))

	var meta, metaError = body.ReadString(0x00)
	if metaError != nil {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}
	meta, _ = strings.CutSuffix(meta, (string)(0x00))

	var blob = gophkeeper.Blob{
		Meta:    meta,
		Content: in.Body,
	}
	rid, storeError := identity.StoreBlob(in.Context(), blob, password)
	if storeError != nil {
		var status = http.StatusInternalServerError
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
	var identity, identityError = e.Repository.Identity(in.Context(), (gophkeeper.Token)(token))
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

	var request struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(in.Body).Decode(&request); err != nil {
		var status = http.StatusBadRequest
		http.Error(out, http.StatusText(status), status)
		return
	}

	var blob, restoreError = identity.RestoreBlob(in.Context(), (gophkeeper.ResourceID)(rid), request.Password)
	if restoreError != nil {
		var status = http.StatusInternalServerError
		if errors.Is(identityError, gophkeeper.ErrBadCredential) {
			status = http.StatusUnauthorized
		}
		http.Error(out, http.StatusText(status), status)
		return
	}

	out.WriteHeader(http.StatusOK)

	var output = bufio.NewWriter(out)
	if _, err := output.WriteString(blob.Meta); err != nil {
		log.Printf("failed to write meta: %s", err.Error())
	}
	if err := output.WriteByte(0x00); err != nil {
		log.Printf("failed to write meta delimiter: %s", err.Error())
	}
	if err := output.Flush(); err != nil {
		log.Printf("failed to flush meta: %s", err.Error())
	}
	if _, err := output.ReadFrom(blob.Content); err != nil {
		log.Printf("failed to write content: %s", err.Error())
	}
}
