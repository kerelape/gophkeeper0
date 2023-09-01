package gophkeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// RestGophkeeper is a remote gophkeeper.
type RestGophkeeper struct {
	Client http.Client
	Server string
}

var _ Gophkeeper = (*RestGophkeeper)(nil)

// Register implements Gophkeeper.
func (g *RestGophkeeper) Register(context context.Context, credential Credential) error {
	var (
		endpoint = fmt.Sprintf("%s/register", g.Server)
		request  struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
	)
	request.Username = credential.Username
	request.Password = credential.Password
	var content, marshalError = json.Marshal(&request)
	if marshalError != nil {
		return marshalError
	}
	response, postError := g.Client.Post(endpoint, "application/json", bytes.NewReader(content))
	if postError != nil {
		return postError
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusConflict:
		return ErrIdentityDuplicate
	}
	return nil
}

// Authenticate implements Gophkeeper.
//
// @todo #28 Implement authentication.
func (*RestGophkeeper) Authenticate(context.Context, Credential) (Token, error) {
	panic("unimplemented")
}

// Identity implements Gophkeeper.
//
// @todo #28 Implement Identity.
func (*RestGophkeeper) Identity(context.Context, Token) (Identity, error) {
	panic("unimplemented")
}
