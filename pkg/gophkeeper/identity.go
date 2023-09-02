package gophkeeper

import (
	"context"
	"errors"
	"io"
)

type (
	// Piece is a piece of encrypted information.
	Piece struct {
		Content []byte // Content of the piece.
		Meta    string // Meta info of the piece.
	}

	// Blob is an encrypted blob.
	Blob struct {
		Content io.ReadCloser // Content of the blob.
		Meta    string        // Meta info of the blob.
	}
)

// ResourceID is id of a resource.
type ResourceID int64

// ResourceType is type of resource stored.
type ResourceType int

const (
	// ResourceTypePiece is a resource type
	// indicating a resource of type Piece.
	ResourceTypePiece ResourceType = iota + 1

	// ResourceTypeBlob is a reource type
	// indicating a resource of type Blob.
	ResourceTypeBlob
)

// Resource is a resource information.
type Resource struct {
	ID   ResourceID
	Type ResourceType
	Meta string
}

var (
	// ErrResourceNotFound is returned when there is no
	// resource with the ResourceID (or it's owned by another identity).
	ErrResourceNotFound = errors.New("resource not found")
)

// Identity is a gophkeeper's identity.
type Identity interface {
	// StorePiece stores a piece and returns its ResourceID.
	StorePiece(ctx context.Context, piece Piece, password string) (ResourceID, error)

	// RestorePiece restores a piece by ResourceID.
	RestorePiece(ctx context.Context, rid ResourceID, password string) (Piece, error)

	// StoreBlob stores a blob and returns its ResourceID.
	StoreBlob(ctx context.Context, blob Blob, password string) (ResourceID, error)

	// RestoreBlob restores a blob by ResourceID.
	RestoreBlob(ctx context.Context, rid ResourceID, password string) (Blob, error)

	// Delete deletes the resource by ResourceID.
	Delete(context.Context, ResourceID) error

	// List returns list of all stored resources.
	List(context.Context) ([]Resource, error)
}
