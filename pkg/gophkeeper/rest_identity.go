package gophkeeper

import (
	"context"
	"net/http"
)

// RestIdentity is rest identity.
type RestIdentity struct {
	Client http.Client
	Server string
	Token  Token
}

var _ Identity = (*RestIdentity)(nil)

// StorePiece implements Identity.
//
// @todo #31 Implement StorePiece on RestIdentity.
func (i *RestIdentity) StorePiece(ctx context.Context, piece Piece, password string) (ResourceID, error) {
	panic("unimplemented")
}

// RestorePiece implements Identity.
//
// @todo #31 Implement RestorePiece on RestIdentity.
func (i *RestIdentity) RestorePiece(ctx context.Context, rid ResourceID, password string) (Piece, error) {
	panic("unimplemented")
}

// StoreBlob implements Identity.
//
// @todo #31 Implement StoreBlob on RestIdentity.
func (i *RestIdentity) StoreBlob(ctx context.Context, blob Blob, password string) (ResourceID, error) {
	panic("unimplemented")
}

// RestoreBlob implements Identity.
//
// @todo #31 Implement RestoreBlob on RestIdentity.
func (i *RestIdentity) RestoreBlob(ctx context.Context, rid ResourceID, password string) (Blob, error) {
	panic("unimplemented")
}

// Delete implements Identity.
//
// @todo #31 Implement Delete on RestIdentity.
func (i *RestIdentity) Delete(context.Context, ResourceID) error {
	panic("unimplemented")
}

// List implements Identity.
//
// @todo #31 Implement List on RestIdentity.
func (i *RestIdentity) List(context.Context) ([]Resource, error) {
	panic("unimplemented")
}
