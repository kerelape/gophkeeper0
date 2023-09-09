package gophkeeper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrIncompatibleAPI is returns when API is not compatible with implementation.
var ErrIncompatibleAPI = errors.New("incompatable API")

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
	case http.StatusCreated:
		return nil
	default:
		return errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
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
		return (Token)(""), errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// Identity implements Gophkeeper.
func (g *RestGophkeeper) Identity(_ context.Context, token Token) (Identity, error) {
	var identity = &RestIdentity{
		Client: g.Client,
		Server: g.Server,
		Token:  token,
	}
	return identity, nil
}
