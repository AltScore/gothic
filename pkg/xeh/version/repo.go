package version

import (
	"context"

	"github.com/looplab/eventhorizon/repo/version"

	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/uuid"
)

type EventStoreReader interface {
	// LoadFrom loads all events from version for the aggregate id from the store.
	LoadFrom(ctx context.Context, id uuid.UUID, version int) ([]eh.Event, error)
}
type MinVersionFunc = func(ctx context.Context, uuid2 uuid.UUID) (int, bool)

// Repo is a middleware that adds version checking to a read repository.
type Repo struct {
	eh.ReadWriteRepo
	eventStore EventStoreReader
}

var _ eh.ReadRepo = (*Repo)(nil)

// NewRepo creates a new Repo.
// Uses the supplied event store to find the min version number of a given stream/aggregate
func NewRepo(repo eh.ReadWriteRepo, eventStore EventStoreReader) *Repo {
	r := &Repo{
		ReadWriteRepo: repo,
		eventStore:    eventStore,
	}
	return r
}

// InnerRepo implements the InnerRepo method of the eventhorizon.ReadRepo interface.
func (r *Repo) InnerRepo(_ context.Context) eh.ReadRepo {
	return r.ReadWriteRepo
}

// IntoRepo tries to convert an eh.ReadRepo into a Repo by recursively looking at
// inner repos. Returns nil if none was found.
func IntoRepo(ctx context.Context, repo eh.ReadRepo) *Repo {
	if repo == nil {
		return nil
	}

	if r, ok := repo.(*Repo); ok {
		return r
	}

	return IntoRepo(ctx, repo.InnerRepo(ctx))
}

// Find implements the Find method of the eventhorizon.ReadModel interface.
// If the context contains a min version set by WithMinVersion it will only
// return an item if its version is at least min version. If a timeout or
// deadline is set on the context it will repeatedly try to get the item until
// either the version matches or the deadline is reached.
func (r *Repo) Find(ctx context.Context, id uuid.UUID) (eh.Entity, error) {
	minVersion, ok := r.findMinVersionNumber(ctx, id)

	if !ok || minVersion < 1 {
		return r.ReadWriteRepo.Find(ctx, id)
	}

	ctx, cancel := version.NewContextWithMinVersionWait(ctx, minVersion)
	defer cancel()

	return r.ReadWriteRepo.Find(ctx, id)
}

// findMinVersionNumber returns the min version number for the given aggregate as given by the event store
func (r *Repo) findMinVersionNumber(ctx context.Context, id uuid.UUID) (int, bool) {
	lastKnown, _ := version.MinVersionFromContext(ctx)

	events, err := r.eventStore.LoadFrom(ctx, id, lastKnown)

	if err != nil {
		return lastKnown, false
	}

	// Find last version used. We assume not specific order of events
	for _, event := range events {
		v := event.Version()
		if v > lastKnown {
			lastKnown = v
		}
	}

	return lastKnown, true
}
