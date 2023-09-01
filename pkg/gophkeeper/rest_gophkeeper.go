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
func (g *RestGophkeeper) Register(ctx context.Context, credential Credential) error {
	var endpoint = fmt.Sprintf("%s/register", g.Server)
	var content, marshalError = json.Marshal(
		map[string]any{
			"username": credential.Username,
			"password": credential.Password,
		},
	)
	if marshalError != nil {
		return marshalError
	}
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodPost, endpoint,
		bytes.NewReader(content),
	)
	if requestError != nil {
		return requestError
	}
	response, postError := g.Client.Do(request)
	if postError != nil {
		return postError
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusConflict:
		return ErrIdentityDuplicate
	case http.StatusOK:
		return nil
	default:
		panic("unexpected response code")
	}
}

// Authenticate implements Gophkeeper.
func (g *RestGophkeeper) Authenticate(ctx context.Context, credential Credential) (Token, error) {
	var endpoint = fmt.Sprintf("%s/login", g.Server)
	var content, marshalError = json.Marshal(
		map[string]any{
			"username": credential.Username,
			"password": credential.Password,
		},
	)
	if marshalError != nil {
		return (Token)(""), marshalError
	}
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodPost, endpoint,
		bytes.NewReader(content),
	)
	if requestError != nil {
		return (Token)(""), requestError
	}
	var response, postError = g.Client.Do(request)
	if postError != nil {
		return (Token)(""), postError
	}
	defer response.Body.Close()
	switch response.StatusCode {
	case http.StatusUnauthorized:
		return (Token)(""), ErrBadCredential
	case http.StatusOK:
		var token = response.Header.Get("Authorization")
		return (Token)(token), nil
	default:
		panic("unexpected response code")
	}
}

// Identity implements Gophkeeper.
//
// @todo #28 Implement Identity.
func (*RestGophkeeper) Identity(context.Context, Token) (Identity, error) {
	panic("unimplemented")
}
