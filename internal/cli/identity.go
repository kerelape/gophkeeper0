package cli

import (
	"context"
	"encoding/json"

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

func (i identity) List(ctx context.Context) ([]resource, error) {
	var resources, resourcesError = i.origin.List(ctx)
	if resourcesError != nil {
		return nil, resourcesError
	}
	var result = make([]resource, 0, len(resources))
	for _, r := range resources {
		var resource resource
		resource.RID = r.ID
		var meta map[string]any
		if err := json.Unmarshal(([]byte)(r.Meta), &meta); err != nil {
			continue
		}
		if value, ok := meta["type"].(int); ok {
			resource.Type = (resourceType)(value)
		} else {
			continue
		}
		if value, ok := meta["description"].(string); ok {
			resource.Description = value
		} else {
			continue
		}
		result = append(result, resource)
	}
	return result, nil
}

func (i identity) StoreCredential(
	ctx context.Context,
	username, password, description string,
	vaultPassword string,
) (gophkeeper.ResourceID, error) {
	var meta, metaError = json.Marshal(
		map[string]any{
			"type":        (int)(resourceTypeCredential),
			"description": description,
		},
	)
	if metaError != nil {
		return -1, metaError
	}
	var content, contentError = json.Marshal(
		map[string]any{
			"username": username,
			"password": password,
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
