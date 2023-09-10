package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"

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

type (
	credentialResource struct {
		description string
		username    string
		password    string
	}
	textResource struct {
		description string
		content     string
	}
	fileResource struct {
		description string
		path        string
	}
	cardResource struct {
		cardInfo
		description string
	}
)

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

func (i identity) StoreText(ctx context.Context, resource textResource, vaultPassword string) (gophkeeper.ResourceID, error) {
	var meta, metaError = json.Marshal(
		map[string]any{
			"type":        (int)(resourceTypeText),
			"description": resource.description,
		},
	)
	if metaError != nil {
		return -1, metaError
	}
	var piece = gophkeeper.Piece{
		Meta:    (string)(meta),
		Content: ([]byte)(resource.content),
	}
	return i.origin.StorePiece(ctx, piece, vaultPassword)
}

func (i identity) RestoreText(ctx context.Context, rid gophkeeper.ResourceID, vaultPassword string) (textResource, error) {
	var piece, pieceError = i.origin.RestorePiece(ctx, rid, vaultPassword)
	if pieceError != nil {
		return textResource{}, pieceError
	}

	var meta struct {
		Type        resourceType `json:"type"`
		Description string       `json:"description"`
	}
	if err := json.Unmarshal(([]byte)(piece.Meta), &meta); err != nil {
		return textResource{}, err
	}
	if meta.Type != resourceTypeText {
		return textResource{}, errors.New("invalid resource type")
	}

	var resource = textResource{
		description: meta.Description,
		content:     (string)(piece.Content),
	}
	return resource, nil
}

func (i identity) StoreFile(ctx context.Context, resource fileResource, vaultPassword string) (gophkeeper.ResourceID, error) {
	var meta, metaError = json.Marshal(
		map[string]any{
			"type":        (int)(resourceTypeFile),
			"description": resource.description,
		},
	)
	if metaError != nil {
		return -1, metaError
	}

	var file, fileError = os.Open(resource.path)
	if fileError != nil {
		return -1, fileError
	}

	var blob = gophkeeper.Blob{
		Meta:    (string)(meta),
		Content: file,
	}
	var rid, ridError = i.origin.StoreBlob(ctx, blob, vaultPassword)
	if ridError != nil {
		return -1, ridError
	}

	return rid, nil
}

func (i identity) RestoreFile(ctx context.Context, rid gophkeeper.ResourceID, path, vaultPassword string) (fileResource, error) {
	var blob, blobError = i.origin.RestoreBlob(ctx, rid, vaultPassword)
	if blobError != nil {
		return fileResource{}, blobError
	}
	defer blob.Content.Close()

	var meta struct {
		Type        resourceType `json:"type"`
		Description string       `json:"description"`
	}
	if err := json.Unmarshal(([]byte)(blob.Meta), &meta); err != nil {
		return fileResource{}, err
	}
	if meta.Type != resourceTypeFile {
		return fileResource{}, errors.New("invalid resource type")
	}

	var file, fileError = os.Create(path)
	if fileError != nil {
		return fileResource{}, fileError
	}
	defer file.Close()

	var input = bufio.NewReader(blob.Content)
	if _, err := input.WriteTo(file); err != nil {
		return fileResource{}, err
	}

	var resource = fileResource{
		path:        file.Name(),
		description: meta.Description,
	}
	return resource, nil
}

func (i identity) StoreCard(ctx context.Context, resource cardResource, vaultPassword string) (gophkeeper.ResourceID, error) {
	var meta, metaError = json.Marshal(
		map[string]any{
			"type":        (int)(resourceTypeCard),
			"description": resource.description,
		},
	)
	if metaError != nil {
		return -1, metaError
	}

	var content, contentError = json.Marshal(
		map[string]any{
			"ccn":    resource.ccn,
			"exp":    resource.exp,
			"cvv":    resource.cvv,
			"holder": resource.holder,
		},
	)
	if contentError != nil {
		return -1, contentError
	}

	var piece = gophkeeper.Piece{
		Meta:    (string)(meta),
		Content: content,
	}
	var rid, ridError = i.origin.StorePiece(ctx, piece, vaultPassword)
	if ridError != nil {
		return -1, ridError
	}

	return rid, nil
}

func (i identity) RestoreCard(ctx context.Context, rid gophkeeper.ResourceID, vaultPassword string) (cardResource, error) {
	var piece, pieceError = i.origin.RestorePiece(ctx, rid, vaultPassword)
	if pieceError != nil {
		return cardResource{}, pieceError
	}

	var meta struct {
		Type        resourceType `json:"type"`
		Description string       `json:"description"`
	}
	if err := json.Unmarshal(([]byte)(piece.Meta), &meta); err != nil {
		return cardResource{}, err
	}
	if meta.Type != resourceTypeCard {
		return cardResource{}, errors.New("invalid resource type")
	}

	var content struct {
		CCN    string `json:"ccn"`
		EXP    string `json:"exp"`
		CVV    string `json:"cvv"`
		Holder string `json:"holder"`
	}
	if err := json.Unmarshal(piece.Content, &content); err != nil {
		return cardResource{}, err
	}

	var resource = cardResource{
		description: meta.Description,
		cardInfo: cardInfo{
			ccn:    content.CCN,
			exp:    content.EXP,
			cvv:    content.CVV,
			holder: content.Holder,
		},
	}
	return resource, nil
}
