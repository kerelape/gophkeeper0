package gophkeeper

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// ErrServerIsDown is returns when server returned an internal server error.
var ErrServerIsDown = errors.New("server is down")

// RestIdentity is rest identity.
type RestIdentity struct {
	Client http.Client
	Server string
	Token  Token
}

var _ Identity = (*RestIdentity)(nil)

// StorePiece implements Identity.
func (i *RestIdentity) StorePiece(ctx context.Context, piece Piece, password string) (ResourceID, error) {
	var endpoint = fmt.Sprintf("%s/vault/piece", i.Server)
	var content, contentError = json.Marshal(
		map[string]any{
			"meta":    piece.Meta,
			"content": base64.RawStdEncoding.EncodeToString(([]byte)(piece.Content)),
		},
	)
	if contentError != nil {
		return -1, contentError
	}
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodPut, endpoint,
		bytes.NewReader(content),
	)
	if requestError != nil {
		return -1, requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))
	request.Header.Set("X-Password", password)

	var response, responseError = i.Client.Do(request)
	if responseError != nil {
		return -1, responseError
	}
	switch response.StatusCode {
	case http.StatusCreated:
		var content struct {
			RID ResourceID `json:"rid"`
		}
		if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
			return -1, errors.Join(
				fmt.Errorf("parse response: %w", err),
				ErrIncompatibleAPI,
			)
		}
		return content.RID, nil
	case http.StatusUnauthorized:
		return -1, ErrBadCredential
	case http.StatusInternalServerError:
		return -1, ErrServerIsDown
	default:
		return -1, errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// RestorePiece implements Identity.
func (i *RestIdentity) RestorePiece(ctx context.Context, rid ResourceID, password string) (Piece, error) {
	var endpoint = fmt.Sprintf("%s/vault/piece/%d", i.Server, rid)
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodGet, endpoint,
		nil,
	)
	if requestError != nil {
		return Piece{}, requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))
	request.Header.Set("X-Password", password)

	var response, responseError = i.Client.Do(request)
	if responseError != nil {
		return Piece{}, responseError
	}
	switch response.StatusCode {
	case http.StatusOK:
		var content = make(map[string]any)
		if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
			return Piece{}, errors.Join(
				fmt.Errorf("parse response: %w", err),
				ErrIncompatibleAPI,
			)
		}
		var piece Piece
		if meta, ok := content["meta"].(string); ok {
			piece.Meta = meta
		} else {
			return Piece{}, errors.Join(
				fmt.Errorf("invalid response body"),
				ErrIncompatibleAPI,
			)
		}
		if content, ok := content["content"].(string); ok {
			var decodedContent, decodedContentError = base64.RawStdEncoding.DecodeString(content)
			if decodedContentError != nil {
				return Piece{}, errors.Join(
					fmt.Errorf("decode content: %w", decodedContentError),
					ErrIncompatibleAPI,
				)
			}
			piece.Content = decodedContent
		} else {
			return Piece{}, errors.Join(
				fmt.Errorf("invalid response body"),
				ErrIncompatibleAPI,
			)
		}
		return piece, nil
	case http.StatusUnauthorized:
		return Piece{}, ErrBadCredential
	case http.StatusInternalServerError:
		return Piece{}, ErrServerIsDown
	default:
		return Piece{}, errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// StoreBlob implements Identity.
func (i *RestIdentity) StoreBlob(ctx context.Context, blob Blob, password string) (ResourceID, error) {
	var endpoint = fmt.Sprintf("%s/vault/blob", i.Server)
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodPut, endpoint,
		blob.Content,
	)
	if requestError != nil {
		return -1, requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))
	request.Header.Set("X-Password", password)
	request.Header.Set("X-Meta", blob.Meta)

	response, responseError := i.Client.Do(request)
	if responseError != nil {
		return -1, responseError
	}

	switch response.StatusCode {
	case http.StatusCreated:
		var content struct {
			RID ResourceID `json:"rid"`
		}
		if err := json.NewDecoder(response.Body).Decode(&content); err != nil {
			return -1, ErrIncompatibleAPI
		}
		return content.RID, nil
	case http.StatusUnauthorized:
		return -1, ErrBadCredential
	case http.StatusInternalServerError:
		return -1, ErrServerIsDown
	default:
		return -1, errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// RestoreBlob implements Identity.
func (i *RestIdentity) RestoreBlob(ctx context.Context, rid ResourceID, password string) (Blob, error) {
	var endpoint = fmt.Sprintf("%s/vault/blob/%d", i.Server, rid)
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodGet, endpoint,
		nil,
	)
	if requestError != nil {
		return Blob{}, requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))
	request.Header.Set("X-Password", password)

	var response, responseError = i.Client.Do(request)
	if responseError != nil {
		return Blob{}, responseError
	}
	switch response.StatusCode {
	case http.StatusOK:
		var blob = Blob{
			Meta:    response.Header.Get("X-Meta"),
			Content: response.Body,
		}
		return blob, nil
	case http.StatusUnauthorized:
		return Blob{}, ErrBadCredential
	case http.StatusInternalServerError:
		return Blob{}, ErrServerIsDown
	default:
		return Blob{}, errors.Join(
			fmt.Errorf("unexpected response status: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// Delete implements Identity.
func (i *RestIdentity) Delete(ctx context.Context, rid ResourceID) error {
	var endpoint = fmt.Sprintf("%s/vault/%d", i.Server, rid)
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodDelete, endpoint,
		nil,
	)
	if requestError != nil {
		return requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))

	response, responseError := i.Client.Do(request)
	if responseError != nil {
		return responseError
	}

	switch response.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusInternalServerError:
		return ErrServerIsDown
	case http.StatusNotFound:
		return ErrResourceNotFound
	default:
		return errors.Join(
			fmt.Errorf("unexpected response code: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}

// List implements Identity.
func (i *RestIdentity) List(ctx context.Context) ([]Resource, error) {
	var endpoint = fmt.Sprintf("%s/vault", i.Server)
	var request, requestError = http.NewRequestWithContext(
		ctx,
		http.MethodGet, endpoint,
		nil,
	)
	if requestError != nil {
		return nil, requestError
	}
	request.Header.Set("Authorization", (string)(i.Token))

	var response, responseError = i.Client.Do(request)
	if responseError != nil {
		return nil, responseError
	}
	defer response.Body.Close()

	switch response.StatusCode {
	case http.StatusOK:
		var responseContent = make(
			[]struct {
				Meta string       `json:"meta"`
				RID  ResourceID   `json:"rid"`
				Type ResourceType `json:"type"`
			},
			0,
		)
		if err := json.NewDecoder(response.Body).Decode(&responseContent); err != nil {
			return nil, err
		}
		var resources = make([]Resource, 0, len(responseContent))
		for _, responseResource := range responseContent {
			resources = append(
				resources,
				Resource{
					ID:   responseResource.RID,
					Type: responseResource.Type,
					Meta: responseResource.Meta,
				},
			)
		}
		return resources, nil
	case http.StatusUnauthorized:
		return nil, ErrBadCredential
	case http.StatusInternalServerError:
		return nil, ErrServerIsDown
	default:
		return nil, errors.Join(
			fmt.Errorf("unexpected response code: %d", response.StatusCode),
			ErrIncompatibleAPI,
		)
	}
}
