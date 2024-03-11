package xeh

import (
	"context"
	"github.com/AltScore/gothic/v2/pkg/xerrors"

	eh "github.com/looplab/eventhorizon"
	localEventBus "github.com/looplab/eventhorizon/eventbus/local"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/uuid"
	"go.uber.org/zap"
)

const (
	RegenerateReadModelsCmdType eh.CommandType = "RegenerateReadModelsCmd"
)

// RegenerateReadModelsCmd is a command that removes a read model and reapply all events to get a refreshed read model.
// This action will collide with current transactions. Execute it with caution.
// The action is idempotent. It can be executed multiple times without side effects.
//
// To regenerate a read model, first register a projector for the aggregate type, then send this command.
type RegenerateReadModelsCmd struct {
	ID   uuid.UUID
	Type eh.AggregateType
}

var _ eh.Command = (*RegenerateReadModelsCmd)(nil)

func (r RegenerateReadModelsCmd) AggregateID() uuid.UUID          { return r.ID }
func (r RegenerateReadModelsCmd) AggregateType() eh.AggregateType { return r.Type }
func (r RegenerateReadModelsCmd) CommandType() eh.CommandType     { return RegenerateReadModelsCmdType }

// ReadModelRegenerator is a command handler that regenerates read models
type ReadModelRegenerator struct {
	logger     *zap.Logger
	eventStore eh.EventStore
	eventBus   eh.EventBus

	readModelRepoByType map[eh.AggregateType]eh.ReadWriteRepo
}

// NewReadModelRegenerator creates a new ReadModelRegenerator that reads old events from the given event store
func NewReadModelRegenerator(eventStore eh.EventStore, logger *zap.Logger) *ReadModelRegenerator {
	return &ReadModelRegenerator{
		logger:              logger,
		eventStore:          eventStore,
		eventBus:            localEventBus.NewEventBus(),
		readModelRepoByType: make(map[eh.AggregateType]eh.ReadWriteRepo),
	}
}

var _ eh.CommandHandler = (*ReadModelRegenerator)(nil)

// Register registers a projector for a given aggregate type to allow the regeneration of its read models.
// Multiple projectors of seam read model can be registered for the same aggregate type.
// Also, multiple aggregate types can be registered.
func (r *ReadModelRegenerator) Register(ctx context.Context, aggType eh.AggregateType, prj *projector.EventHandler, repo eh.ReadWriteRepo) error {
	r.readModelRepoByType[aggType] = repo

	return r.eventBus.AddHandler(ctx, eh.MatchAggregates{aggType}, prj)
}

// Regenerate removes the read model and reapply all events to get a refreshed read model.
func (r *ReadModelRegenerator) Regenerate(ctx context.Context, aggType eh.AggregateType, id uuid.UUID) error {
	cmd := RegenerateReadModelsCmd{
		ID:   id,
		Type: aggType,
	}

	return r.HandleCommand(ctx, cmd)
}

func (r *ReadModelRegenerator) HandleCommand(ctx context.Context, command eh.Command) error {
	cmd, ok := command.(RegenerateReadModelsCmd)
	if !ok {
		return eh.ErrCommandNotRegistered
	}

	r.logger.Info("Regenerating read models", zap.String("agg_type", string(cmd.Type)), zap.String("agg_id", cmd.ID.String()))

	if err := r.removeReadModels(ctx, cmd.Type, cmd.ID); err != nil {
		return err
	}

	return r.replayEvents(ctx, cmd.Type, cmd.ID)
}

// removeReadModels removes read models for the given aggregate id
func (r *ReadModelRegenerator) removeReadModels(ctx context.Context, aggregateType eh.AggregateType, id uuid.UUID) error {
	repo, found := r.readModelRepoByType[aggregateType]
	if !found {
		return xerrors.NewNotFoundError("read mode repo", "not found for type %s", aggregateType)

	}
	err := repo.Remove(ctx, id)
	if err == nil || IsEHNotFound(err) {
		return nil
	}
	return err
}

// replayEvents replays events for the given aggregate id from the first event
// Events are not published outside the own event bus
func (r *ReadModelRegenerator) replayEvents(ctx context.Context, aggregateType eh.AggregateType, id uuid.UUID) error {
	events, err := r.eventStore.Load(ctx, id)
	if err != nil {
		return err
	}

	if len(events) == 0 {
		return nil
	}

	if events[0].AggregateType() != aggregateType {
		return eh.ErrMismatchedEventAggregateTypes
	}

	for _, event := range events {
		if err := r.eventBus.HandleEvent(ctx, event); err != nil {
			return err
		}
	}

	return nil
}
