package entity

import "context"

// Updatable is an interface for entities that can update their metadata before saving to repository
type Updatable interface {
	// UpdateMetadata updates the metadata of the entity before saving to repository. Returns the updated metadata.
	UpdateMetadata(ctx context.Context)
}
