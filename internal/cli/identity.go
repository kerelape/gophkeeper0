package cli

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kerelape/gophkeeper/pkg/gophkeeper"
)

type identity struct {
	origin gophkeeper.Identity
}

type resourceType int

const (
	resourceTypeCredential resourceType = iota
	resourceTypeText
	resourceTypeFile
	resourceTypeCard
)

func (r resourceType) String() string {
	switch r {
	case resourceTypeCredential:
		return "Credential"
	case resourceTypeText:
		return "Text"
	case resourceTypeFile:
		return "File"
	case resourceTypeCard:
		return "Card"
	default:
		panic("unknown resource type")
	}
}

type resource struct {
	RID         gophkeeper.ResourceID
	Description string
	Type        resourceType
}

type credentialResource struct {
	description string
	username    string
	password    string
}

func (i identity) List(ctx context.Context) ([]resource, error) {
	var resources, resourcesError = i.origin.List(ctx)
	if resourcesError != nil {
		return nil, resourcesError
	}
	var result = make([]resource, 0, len(resources))
	for _, r := range resources {
		var resource resource
		resource.RID = r.ID
		var meta struct {
			Type        resourceType `json:"type"`
			Description string       `json:"description"`
		}
		if err := json.Unmarshal(([]byte)(r.Meta), &meta); err != nil {
			continue
		}
		resource.Type = meta.Type
		resource.Description = meta.Description
		result = append(result, resource)
	}
	return result, nil
}

func (i identity) StoreCredential(ctx context.Context, cred credentialResource, vaultPassword string) (gophkeeper.ResourceID, error) {
	var meta, metaError = json.Marshal(
		map[string]any{
			"type":        (int)(resourceTypeCredential),
			"description": cred.description,
		},
	)
	if metaError != nil {
		return -1, metaError
	}
	var content, contentError = json.Marshal(
		map[string]any{
			"username": cred.username,
			"password": cred.password,
		},
	)
	if contentError != nil {
		return -1, contentError
	}
	var piece = gophkeeper.Piece{
		Meta:    (string)(meta),
		Content: content,
	}
	return i.origin.StorePiece(ctx, piece, vaultPassword)
}

func (i identity) RestoreCredential(ctx context.Context, rid gophkeeper.ResourceID, vaultPassword string) (credentialResource, error) {
	var piece, pieceError = i.origin.RestorePiece(ctx, rid, vaultPassword)
	if pieceError != nil {
		return credentialResource{}, pieceError
	}

	var meta struct {
		Type        resourceType `json:"type"`
		Description string       `json:"description"`
	}
	if err := json.Unmarshal(([]byte)(piece.Meta), &meta); err != nil {
		return credentialResource{}, err
	}
	if meta.Type != resourceTypeCredential {
		return credentialResource{}, errors.New("invalid resource type")
	}

	var content struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.Unmarshal(piece.Content, &content); err != nil {
		return credentialResource{}, err
	}

	var res = credentialResource{
		description: meta.Description,
		username:    content.Username,
		password:    content.Password,
	}
	return res, nil
}
