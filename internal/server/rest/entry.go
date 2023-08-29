package rest

import "net/http"

// Entry is the REST api entry.
//
// @todo #7 Implement authentication.
// @todo #7 Implement registration.
// @todo #7 Implement text storage.
// @todo #7 Implement blob storage.
// @todo #7 Implement bank card storage.
type Entry struct {
}

// Route routes Entry into an http.Handler.
func (e *Entry) Route() http.Handler {
	return nil
}
