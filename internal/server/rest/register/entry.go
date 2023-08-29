package register

import "net/http"

// Entry is login entry.
type Entry struct {
}

// Route routes login entry.
func (e *Entry) Route() http.Handler {
	return nil
}
